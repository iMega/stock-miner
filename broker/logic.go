package broker

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/uuid"
	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
)

func (b *Broker) run() {
	pricerCh := make(chan domain.Message)

	outCh := make(chan domain.PriceReceiptMessage)
	psCh := make(chan domain.PriceReceiptMessage)
	sellCh := make(chan domain.Slot)

	operationCh := make(chan domain.Message)

	confirmBuyCh := make(chan domain.Message)
	b.confirmBuyWorker(confirmBuyCh, operationCh)

	confirmSellCh := make(chan domain.Message)
	b.confirmSellWorker(confirmSellCh, operationCh)

	b.queueOperation(operationCh, confirmBuyCh, confirmSellCh)

	w1 := b.pricerWorker(pricerCh, outCh, psCh)
	w2 := b.makePriceStorageChannel(psCh)
	b.solveWorker(outCh, sellCh, confirmBuyCh)
	b.sellWorker(sellCh, confirmSellCh)

	delay := cron.DelayIfStillRunning(&logger{log: b.logger})
	b.cron.AddJob("@every 2s", delay(cron.FuncJob(func() {
		if w1.WaitingQueueSize()+w2.WaitingQueueSize() > 20 {
			b.logger.Debugf("WaitingQueueSize = %d", w1.WaitingQueueSize()+w2.WaitingQueueSize())
		}

		b.StockStorage.StockItemApprovedAll(context.Background(), pricerCh)
	})))
}

func (b *Broker) pricerWorker(in chan domain.Message, out, out2 chan domain.PriceReceiptMessage) *workerpool.WorkerPool {
	wp := workerpool.New(5)

	go func() {
		for m := range in {
			msg := m
			wp.Submit(func() {
				if msg.Error != nil {
					b.logger.Errorf("message has error, %W", msg.Error)

					return
				}

				res, err := b.getPrice(msg)
				if err != nil {
					b.logger.Errorf("failed getting price, %s", err)

					return
				}

				if !b.SMAStack.Add(res.Transaction.Slot.Ticker, res.Price) {
					return
				}

				r := domain.PriceReceiptMessage{
					Email:     res.Transaction.Slot.Email,
					Price:     res.Price,
					StockItem: res.Transaction.Slot.StockItem,
				}

				out <- r
				out2 <- r
			})
		}
	}()

	return wp
}

func (b *Broker) makePriceStorageChannel(in chan domain.PriceReceiptMessage) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for task := range in {
			t := task
			wp.Submit(func() {
				err := b.StockStorage.AddMarketPrice(context.Background(), t)
				if err != nil {
					b.logger.Errorf("failed to add market price, %s", err)
				}
			})
		}
	}()

	return wp
}

func (b *Broker) solveWorker(
	in chan domain.PriceReceiptMessage,
	sellCh chan domain.Slot,
	confirmBuyCh chan domain.Message,
) *workerpool.WorkerPool {
	wp := workerpool.New(5)

	go func() {
		for task := range in {
			t := task
			wp.Submit(func() {
				msg := domain.Message{
					Price: t.Price,
					Transaction: domain.Transaction{
						Slot: domain.Slot{
							Email:     t.Email,
							StockItem: t.StockItem,
						},
					},
				}
				err := b.solver(msg, sellCh, confirmBuyCh)
				if err != nil {
					b.logger.Errorf("failed to solve, %s", err)
				}
			})
		}
	}()

	return wp
}

func (b *Broker) solver(
	msg domain.Message,
	sellCh chan domain.Slot,
	confirmBuyCh chan domain.Message,
) error {
	frame, err := b.SMAStack.Get(msg.Transaction.Slot.StockItem.Ticker)
	if err != nil {
		return fmt.Errorf("failed getting frame from stack, %s", err)
	}

	if !frame.IsFull() {
		return fmt.Errorf("frame is not full") //nil
	}

	ctx, err := b.contextWithCreds(context.Background(), msg.Transaction.Slot.Email)
	if err != nil {
		return fmt.Errorf("failed getting creds, %w", err)
	}

	settings, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting settings, %w", err)
	}

	slots, err := b.Stack.Slot(ctx, msg.Transaction.Slot.StockItem.FIGI)
	if err != nil {
		return fmt.Errorf("failed getting slot, %s", err)

	}

	sellSlots := getItemsForSale(slots, frame.Last())
	for _, slot := range sellSlots {
		sellCh <- slot
	}

	if settings.Slot.Volume <= len(slots) {
		return nil
	}

	minPrice := minBuyingPrice(slots)
	if minPrice == 0 {
		return nil
	}

	if minPrice-settings.Slot.ModificatorMinPrice >= msg.Price {
		return nil
	}

	isTrendUp, err := b.SMAStack.IsTrendUp(msg.Transaction.Slot.StockItem.Ticker)
	if err != nil {
		return fmt.Errorf("failed getting trend, %s", err)
	}

	if isTrendUp {
		return nil
	}

	// buy
	emptyTr := domain.Transaction{
		Slot: domain.Slot{
			ID:          uuid.NewID().String(),
			Email:       msg.Transaction.Email,
			StockItem:   msg.Transaction.Slot.StockItem,
			SlotID:      len(slots) + 1,
			StartPrice:  frame.Prev(),
			ChangePrice: frame.Last(),
			Qty:         1,
		},
		BuyAt: time.Now(),
	}

	tr, err := b.buy(ctx, emptyTr)
	if err != nil {
		return fmt.Errorf("failed to buy stock item, %s", err)
	}

	confirmBuyCh <- domain.Message{Transaction: tr}

	return nil
}

func (b *Broker) buyWorker(in chan domain.Slot) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for task := range in {
			t := task
			wp.Submit(func() {
				_ = t
			})
		}
	}()

	return wp
}

func (b *Broker) sellWorker(
	in chan domain.Slot,
	confirmSellCh chan domain.Message,
) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for task := range in {
			t := task
			wp.Submit(func() {
				ctx, err := b.contextWithCreds(context.Background(), t.Email)
				if err != nil {
					b.logger.Errorf("failed getting creds, %s", err)

					return
				}

				tr, err := b.StockStorage.Transaction(ctx, t.ID)
				if err != nil {
					b.logger.Errorf("failed getting transaction, %s", err)

					return
				}

				upTr, err := b.sell(ctx, tr)
				if err != nil {
					b.logger.Errorf("failed to sell items, %s", err)

					return
				}

				confirmSellCh <- domain.Message{
					Transaction: upTr,
				}
			})
		}
	}()

	return wp
}

func (b *Broker) confirmSellWorker(confirmSellCh, operationCh chan domain.Message) *workerpool.WorkerPool {
	wp := workerpool.New(100)

	go func() {
		for m := range confirmSellCh {
			msg := m

			wp.Submit(func() {
				if err := b.confirmSellJob(msg); err != nil {
					b.logger.Errorf("failed to confirm sell, %s", err)
					msg.RetryCount++
					operationCh <- msg
				}
			})
		}
	}()

	return wp
}

func (b *Broker) getPrice(msg domain.Message) (domain.Message, error) {
	result := msg

	ctx, err := b.contextWithCreds(
		context.Background(),
		msg.Transaction.Slot.Email,
	)
	if err != nil {
		return result, fmt.Errorf("failed getting creds, %W", err)
	}

	ob, err := b.Market.OrderBook(ctx, msg.Transaction.Slot.StockItem)
	if err != nil {
		return result, err
	}

	result.Price = ob.LastPrice

	return result, nil
}

func getItemsForSale(slots []domain.Slot, price float64) []domain.Slot {
	result := []domain.Slot{}
	p := decimal.NewFromFloat(price)

	for _, slot := range slots {
		if slot.BuyingPrice == 0 {
			continue
		}

		if decimal.NewFromFloat(slot.TargetPrice).LessThanOrEqual(p) {
			result = append(result, slot)
		}
	}

	return result
}

func (b *Broker) queueOperation(
	in, confirmBuyCh, confirmSellCh chan domain.Message,
) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for m := range in {
			msg := m
			<-time.After(4 * time.Second)

			wp.Submit(func() {
				newMsg, op, err := processOperation(msg)
				if err != nil {
					b.logger.Errorf("failed to process operation, %w", err)

					return
				}

				if op == domain.BUY {
					confirmBuyCh <- newMsg
				}

				if op == domain.SELL {
					confirmSellCh <- newMsg
				}
			})

		}
	}()

	return wp
}

func processOperation(msg domain.Message) (domain.Message, domain.OperationType, error) {
	msg.RetryCount++
	if msg.RetryCount > 60 {
		return msg, "", fmt.Errorf(
			"the maximum number of attempts to receive the operation has been reached, id:%s",
			msg.Transaction.ID,
		)
	}

	if msg.Transaction.BuyingPrice == 0 {
		return msg, domain.BUY, nil
	}

	if msg.Transaction.SalePrice == 0 {
		return msg, domain.SELL, nil
	}

	return msg, "", errors.New("unknown transaction type")
}

func (b *Broker) confirmBuyWorker(
	confirmBuyCh chan domain.Message,
	operationCh chan domain.Message,
) *workerpool.WorkerPool {
	wp := workerpool.New(100)

	go func() {
		for m := range confirmBuyCh {
			msg := m

			wp.Submit(func() {
				if err := b.confirmBuyJob(msg.Transaction); err != nil {
					b.logger.Errorf("failed to confirm buy, %s", err)
					msg.RetryCount++
					operationCh <- msg
				}
			})
		}
	}()

	return wp
}

func filterOperationByOrderID(trs []domain.Transaction, orderID string) (domain.Transaction, error) {
	for _, t := range trs {
		if t.BuyOrderID == orderID {
			return t, nil
		}
	}

	return domain.Transaction{}, fmt.Errorf("operation does not exist")
}

func (b *Broker) contextWithCreds(ctxIn context.Context, email string) (context.Context, error) {
	ctx := contexkey.WithEmail(ctxIn, email)

	settings, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		b.logger.Errorf("failed getting settings, %s", err)

		return nil, err
	}

	cred := settings.MarketCredentials[settings.MarketProvider]
	ctx = contexkey.WithToken(ctx, cred.Token)
	ctx = contexkey.WithAPIURL(ctx, cred.APIURL)

	return ctx, nil
}

func minBuyingPrice(slots []domain.Slot) float64 {
	if len(slots) == 0 {
		return -1
	}

	var byuing []float64
	for _, slot := range slots {
		byuing = append(byuing, slot.BuyingPrice)
	}

	sort.Float64s(byuing)

	return byuing[0]
}
