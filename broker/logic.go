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
	"github.com/imega/stock-miner/worker"
	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
)

const (
	maxQueues  = 20
	five       = 5
	hundred    = 100
	delayTask  = 4 * time.Second
	retryCount = 60
)

var (
	errMaxAttempts        = errors.New("the maximum number of attempts to receive the operation has been reached")
	errUnknownTransaction = errors.New("unknown transaction type")
	errFrameNotFull       = errors.New("frame is not full")
	errRangeIsZero        = errors.New("range is zero")
	errOperationNotExist  = errors.New("operation does not exist")
)

func (b *Broker) run() {
	b.confirmBuyWorker(b.Traffic.ConfirmBuyCh, b.Traffic.OperationCh)
	b.confirmSellWorker(b.Traffic.ConfirmSellCh, b.Traffic.OperationCh)
	b.queueOperation(b.Traffic.OperationCh, b.Traffic.ConfirmBuyCh, b.Traffic.ConfirmSellCh)
	b.sellWorker(b.Traffic.SellCh, b.Traffic.ConfirmSellCh)

	delay := cron.DelayIfStillRunning(&logger{log: b.logger})

	timeLayout := "15:04"
	start, _ := time.Parse(timeLayout, "00:00")
	end, _ := time.Parse(timeLayout, "00:05")

	_, err := b.cron.AddJob("@every 2s", delay(cron.FuncJob(func() {
		wd := time.Now().Weekday()
		if !b.isDevMode && (wd == time.Sunday || wd == time.Saturday) {
			return
		}

		if inTimeSpan(start, end, time.Now()) {
			b.SMAStack.Reset()
		}

		b.StockStorage.StockItemApprovedAll(context.Background(), b.Traffic.ApprovedCh)
	})))
	if err != nil {
		b.logger.Errorf("failed to add job to cron, %w", err)
	}
}

func (this *Broker) stockItemApprovedWorker(w worker.Worker) {
	for m := range this.Traffic.ApprovedCh {
		msg := m

		w.Submit(func() {
			if msg.Error != nil {
				this.logger.Errorf(
					"pricer worker reports about the message has error, %w",
					msg.Error,
				)

				return
			}

			res, err := this.getPrice(msg)
			if err != nil {
				this.logger.Errorf("failed getting price, %s", err)

				return
			}

			if !this.SMAStack.Add(res.Transaction.Slot.Ticker, res.Price) {
				return
			}

			r := domain.PriceReceiptMessage{
				Email:     res.Transaction.Slot.Email,
				Price:     res.Price,
				StockItem: res.Transaction.Slot.StockItem,
			}

			this.Traffic.PriceReceiptCh <- r
			this.Traffic.PriceReceiptStoreCh <- r
		})
	}
}

func (b *Broker) makePriceStorageWorker(w worker.Worker) {
	for task := range b.Traffic.PriceReceiptStoreCh {
		t := task

		w.Submit(func() {
			err := b.StockStorage.AddMarketPrice(context.Background(), t)
			if err != nil {
				b.logger.Errorf("failed to add market price, %s", err)
			}
		})
	}
}

func (b *Broker) priceReceiptWorker(w worker.Worker) {
	for t := range b.Traffic.PriceReceiptCh {
		task := t

		w.Submit(func() {
			msg := domain.Message{
				Price: task.Price,
				Transaction: domain.Transaction{
					Slot: domain.Slot{
						Email:     task.Email,
						StockItem: task.StockItem,
					},
				},
			}

			if err := b.branchBuyOrSell(msg); err != nil {
				b.logger.Errorf("solve worker reports, %s", err)

				if errors.Is(err, errRangeIsZero) {
					ctx := context.Background()
					r, rErr := b.Pricer.Range(ctx, task.StockItem)
					if rErr != nil {
						b.logger.Errorf(
							"failed getting stock item range, %s",
							rErr,
						)
					}

					frame, err := b.SMAStack.Get(task.StockItem.Ticker)
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
}

func priceInRange(frame domain.SMAFrame, p float64) bool {
	h, l := frame.RangeHL()

	high := decimal.NewFromFloat(h)
	low := decimal.NewFromFloat(l)
	price := decimal.NewFromFloat(p)

	return price.LessThan(high) && price.GreaterThan(low)
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
					b.logger.Errorf("failed to confirm sell job, %s", err)
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

func getItemsForSale(slots []domain.Slot, price, prevPrice float64) []domain.Slot {
	result := []domain.Slot{}
	p := decimal.NewFromFloat(price)
	pp := decimal.NewFromFloat(prevPrice)

	for _, slot := range slots {
		if slot.BuyingPrice == 0 {
			continue
		}

		target := decimal.NewFromFloat(slot.TargetPrice)
		if target.LessThanOrEqual(p) && target.LessThanOrEqual(pp) {
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

func filterOperationByOrderID(trs []domain.Transaction, orderID string) (domain.Transaction, error) {
	for _, t := range trs {
		if t.BuyOrderID == orderID {
			return t, nil
		}
	}

	return domain.Transaction{},
		fmt.Errorf("%w, orderID=%s", errOperationNotExist, orderID)
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

func minBuyingPrice(slots []domain.Slot, buyingPrice float64) float64 {
	if buyingPrice > 0 {
		slots = append(slots, domain.Slot{BuyingPrice: buyingPrice})
	}

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

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}

	if start.Equal(end) {
		return check.Equal(start)
	}

	return !start.After(check) || !end.Before(check)
}
