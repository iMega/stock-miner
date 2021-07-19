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

const maxQueues = 20

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
	b.solveWorker(solveWorkerInput{
		MessageCh:    outCh,
		SellCh:       sellCh,
		ConfirmBuyCh: confirmBuyCh,
	})
	b.sellWorker(sellCh, confirmSellCh)

	delay := cron.DelayIfStillRunning(&logger{log: b.logger})

	_, err := b.cron.AddJob("@every 2s", delay(cron.FuncJob(func() {
		if w1.WaitingQueueSize()+w2.WaitingQueueSize() > maxQueues {
			b.logger.Debugf("WaitingQueueSize = %d", w1.WaitingQueueSize()+w2.WaitingQueueSize())
		}

		b.StockStorage.StockItemApprovedAll(context.Background(), pricerCh)
	})))
	if err != nil {
		b.logger.Errorf("failed to add jiob to cron, %w", err)
	}
}

const (
	five    = 5
	hundred = 100
)

func (b *Broker) pricerWorker(
	in chan domain.Message,
	out, out2 chan domain.PriceReceiptMessage,
) *workerpool.WorkerPool {
	wp := workerpool.New(five)

	go func() {
		for m := range in {
			msg := m

			wp.Submit(func() {
				if msg.Error != nil {
					b.logger.Errorf(
						"pricer worker reports about the message has error, %w",
						msg.Error,
					)

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

type solveWorkerInput struct {
	MessageCh    chan domain.PriceReceiptMessage
	SellCh       chan domain.Slot
	ConfirmBuyCh chan domain.Message
	RequestRange chan domain.Slot
}

func (b *Broker) solveWorker(in solveWorkerInput) *workerpool.WorkerPool {
	wp := workerpool.New(five)

	go func() {
		for task := range in.MessageCh {
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

				input := solverInput{
					Message:      msg,
					SellCh:       in.SellCh,
					ConfirmBuyCh: in.ConfirmBuyCh,
				}
				if err := b.solver(input); err != nil {
					b.logger.Errorf("solve worker reports, %s", err)

					if err == errRangeIsZero {
						ctx := context.Background()
						r, rErr := b.Pricer.Range(ctx, t.StockItem)
						if rErr != nil {
							b.logger.Errorf(
								"failed getting stock item range, %s",
								rErr,
							)
						}

						frame, err := b.SMAStack.Get(t.StockItem.Ticker)
						if err != nil {
							b.logger.Errorf(
								"failed getting frame from stack, %s",
								err,
							)
						}

						frame.SetRangeHL(r.High, r.Low)
					}
				}
			})
		}
	}()

	return wp
}

var (
	errFrameNotFull = errors.New("frame is not full")
	errRangeIsZero  = errors.New("range is zero")
)

type solverInput struct {
	Message      domain.Message
	SellCh       chan domain.Slot
	ConfirmBuyCh chan domain.Message
}

func (b *Broker) solver(in solverInput) error {
	frame, err := b.SMAStack.Get(in.Message.Transaction.Slot.StockItem.Ticker)
	if err != nil {
		return fmt.Errorf("failed getting frame from stack, %w", err)
	}

	h, l := frame.RangeHL()
	if h == 0 || l == 0 {
		return errRangeIsZero
	}

	if !frame.IsFull() {
		return errFrameNotFull
	}

	ctx, err := b.contextWithCreds(
		context.Background(),
		in.Message.Transaction.Slot.Email,
	)
	if err != nil {
		return fmt.Errorf("failed getting creds, %w", err)
	}

	settings, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting settings, %w", err)
	}

	slots, err := b.Stack.Slot(ctx, in.Message.Transaction.Slot.StockItem.FIGI)
	if err != nil {
		return fmt.Errorf("failed getting slot, %w", err)
	}

	sellSlots := getItemsForSale(slots, frame.Last())
	for _, slot := range sellSlots {
		in.SellCh <- slot
	}

	if settings.Slot.Volume <= len(slots) {
		return nil
	}

	minPrice := minBuyingPrice(slots)
	if minPrice == 0 {
		return nil
	}

	if minPrice-settings.Slot.ModificatorMinPrice >= in.Message.Price {
		return nil
	}

	isTrendUp, err := b.SMAStack.IsTrendUp(in.Message.Transaction.Slot.StockItem.Ticker)
	if err != nil {
		return fmt.Errorf("failed getting trend, %w", err)
	}

	if isTrendUp {
		return nil
	}

	// buy
	emptyTr := domain.Transaction{
		Slot: domain.Slot{
			ID:          uuid.NewID().String(),
			Email:       in.Message.Transaction.Email,
			StockItem:   in.Message.Transaction.Slot.StockItem,
			SlotID:      len(slots) + 1,
			StartPrice:  frame.Prev(),
			ChangePrice: frame.Last(),
			Qty:         1,
		},
		BuyAt: time.Now(),
	}

	tr, err := b.buy(ctx, emptyTr)
	if err != nil {
		return fmt.Errorf("failed to buy stock item, %w", err)
	}

	in.ConfirmBuyCh <- domain.Message{Transaction: tr}

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
	wp := workerpool.New(hundred)

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
		return result, fmt.Errorf("failed getting creds, %w", err)
	}

	ob, err := b.Market.OrderBook(ctx, msg.Transaction.Slot.StockItem)
	if err != nil {
		return result, fmt.Errorf("failed getting orderbook, %w", err)
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

const delayTask = 4 * time.Second

func (b *Broker) queueOperation(
	in, confirmBuyCh, confirmSellCh chan domain.Message,
) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for m := range in {
			msg := m

			<-time.After(delayTask)

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

var (
	errMaxAttempts        = errors.New("the maximum number of attempts to receive the operation has been reached")
	errUnknownTransaction = errors.New("unknown transaction type")
)

const retryCount = 60

func processOperation(msg domain.Message) (domain.Message, domain.OperationType, error) {
	msg.RetryCount++
	if msg.RetryCount > retryCount {
		return msg,
			"",
			fmt.Errorf("%w, id:%s", errMaxAttempts, msg.Transaction.ID)
	}

	if msg.Transaction.BuyingPrice == 0 {
		return msg, domain.BUY, nil
	}

	if msg.Transaction.SalePrice == 0 {
		return msg, domain.SELL, nil
	}

	return msg, "", errUnknownTransaction
}

func (b *Broker) confirmBuyWorker(
	confirmBuyCh chan domain.Message,
	operationCh chan domain.Message,
) *workerpool.WorkerPool {
	wp := workerpool.New(hundred)

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

var errOperationNotExist = errors.New("operation does not exist")

func filterOperationByOrderID(trs []domain.Transaction, orderID string) (domain.Transaction, error) {
	for _, t := range trs {
		if t.BuyOrderID == orderID {
			return t, nil
		}
	}

	return domain.Transaction{}, errOperationNotExist
}

func (b *Broker) contextWithCreds(ctxIn context.Context, email string) (context.Context, error) {
	ctx := contexkey.WithEmail(ctxIn, email)

	settings, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		b.logger.Errorf("failed getting settings, %s", err)

		return nil, fmt.Errorf("failed getting settings, %w", err)
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

	byuing := make([]float64, len(slots))
	for i, slot := range slots {
		byuing[i] = slot.BuyingPrice
	}

	sort.Float64s(byuing)

	return byuing[0]
}
