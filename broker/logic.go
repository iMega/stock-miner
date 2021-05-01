package broker

import (
	"context"
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
	inCh := make(chan domain.PriceReceiptMessage)
	outCh := make(chan domain.PriceReceiptMessage)
	psCh := make(chan domain.PriceReceiptMessage)
	sellCh := make(chan domain.Slot)
	w1 := b.makePricerChannel(inCh, outCh, psCh)
	w2 := b.makePriceStorageChannel(psCh)
	b.noName(outCh, sellCh)
	b.sellWorker(sellCh)

	delay := cron.DelayIfStillRunning(&logger{log: b.logger})
	b.cron.AddJob("@every 2s", delay(cron.FuncJob(func() {
		if w1.WaitingQueueSize()+w2.WaitingQueueSize() > 20 {
			b.logger.Info("STOPPED")
			return
		}
		b.StockStorage.StockItemApprovedAll(context.Background(), inCh)
	})))
}

func (b *Broker) makePricerChannel(in, out, out2 chan domain.PriceReceiptMessage) *workerpool.WorkerPool {
	wp := workerpool.New(5)

	go func() {
		for task := range in {
			t := task
			wp.Submit(func() {
				// res, err := b.Pricer.GetPrice(context.Background(), t)
				// if err != nil {
				// 	b.logger.Errorf("failed getting price from YF, %s", err)
				// }
				// res.Error = err

				// b.logger.Infof("-------- %s %f", res.Ticker, res.Price)

				res, err := b.getPrice(t)
				if err != nil {
					b.logger.Errorf("failed getting price, %s", err)
					return
				}

				if !b.SMAStack.Add(res.Ticker, res.Price) {
					return
				}

				out <- res
				out2 <- res
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

func (b *Broker) noName(in chan domain.PriceReceiptMessage, sellCh chan domain.Slot) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for task := range in {
			t := task
			wp.Submit(func() {
				frame, err := b.SMAStack.Get(t.Ticker)
				if err != nil {
					b.logger.Errorf("failed getting frame from stack, %s", err)
					return
				}

				if !frame.IsFull() {
					b.logger.Error("frame is not full")
					return
				}

				ctx := contexkey.WithEmail(context.Background(), t.Email)
				settings, err := b.SettingsStorage.Settings(ctx)
				if err != nil {
					b.logger.Errorf("failed getting settings, %s", err)
					return
				}

				slots, err := b.Stack.Slot(ctx, t.FIGI)
				if err != nil {
					b.logger.Errorf("failed getting slot, %s", err)
					return
				}

				sellSlots := getItemsForSale(slots, frame.Last)
				for _, slot := range sellSlots {
					sellCh <- slot
				}

				if settings.Slot.Volume <= len(slots) {
					return
				}

				var byuing []float64
				for _, slot := range slots {
					byuing = append(byuing, slot.BuyingPrice)
				}
				sort.Float64s(byuing)

				trend, err := b.SMAStack.IsTrendUp(t.Ticker)
				if err != nil {
					b.logger.Errorf("failed getting trend, %s", err)
					return
				}

				if trend || len(byuing) > 0 && byuing[0]-settings.Slot.ModificatorMinPrice >= t.Price {
					return
				}
				//buy
				cred := settings.MarketCredentials[settings.MarketProvider]
				ctx = contexkey.WithToken(ctx, cred.Token)
				ctx = contexkey.WithAPIURL(ctx, cred.APIURL)

				emptyTr := domain.Transaction{
					Slot: domain.Slot{
						ID:          uuid.NewID().String(),
						Email:       t.Email,
						StockItem:   t.StockItem,
						SlotID:      len(slots) + 1,
						StartPrice:  frame.Prev(),
						ChangePrice: frame.Last,
						Qty:         1,
					},
					BuyAt: time.Now(),
				}

				tr, err := b.buy(ctx, emptyTr)
				if err != nil {
					b.logger.Errorf("failed to buy item, %s", err)
					return
				}

				trs, err := b.Market.Operations(
					ctx,
					domain.OperationInput{
						From:          tr.BuyAt,
						To:            tr.BuyAt.Add(time.Minute),
						OperationType: "Buy",
					},
				)
				if err != nil {
					b.logger.Error(err)
					return
				}

				filteredTR, err := filterOperationByOrderID(trs, tr.BuyOrderID)
				if err != nil {
					b.logger.Error(err)
					return
				}

				tr.BuyingPrice = filteredTR.BuyingPrice
				tr.Slot.AmountSpent = filteredTR.Slot.AmountSpent
				tr.TargetPrice = calcTargetPrice(
					settings.MarketCommission,
					tr.BuyingPrice,
					settings.GrossMargin,
				)

				profit, _ := decimal.NewFromFloat(tr.TargetPrice).Sub(decimal.NewFromFloat(tr.BuyingPrice)).Float64()
				tr.Profit = profit

				if err := b.confirmBuy(ctx, tr); err != nil {
					b.logger.Errorf("failed to confirm transaction, %s", err)
					return
				}
			})
		}
	}()

	return wp
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

func (b *Broker) sellWorker(in chan domain.Slot) *workerpool.WorkerPool {
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

func filterOperationByOrderID(trs []domain.Transaction, orderID string) (domain.Transaction, error) {
	for _, t := range trs {
		if t.BuyOrderID == orderID {
			return t, nil
		}
	}

	return domain.Transaction{}, fmt.Errorf("operation does not exist")
}

func (b *Broker) getPrice(msg domain.PriceReceiptMessage) (domain.PriceReceiptMessage, error) {
	result := msg

	ctx := contexkey.WithEmail(context.Background(), msg.Email)
	settings, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		return result, err
	}

	cred := settings.MarketCredentials[settings.MarketProvider]
	ctx = contexkey.WithToken(ctx, cred.Token)
	ctx = contexkey.WithAPIURL(ctx, cred.APIURL)
	ob, err := b.Market.OrderBook(ctx, msg.StockItem)
	if err != nil {
		return result, err
	}

	result.Price = ob.LastPrice

	return result, nil
}

// формула расчета целевой цены для продажи
//
// ценаПокупки+(ценаПокупки/100*комиссия) = затраты
// затраты + (затраты / 100 * маржа%) = ЦенаПродажиБезКомиссии
// ЦенаПродажиБезКомиссии+(ЦенаПродажиБезКомиссии/100*комиссия) = ЦенаПродажи
func calcTargetPrice(commission, buyingPrice, margin float64) float64 {
	c := decimal.NewFromFloat(commission)
	bp := decimal.NewFromFloat(buyingPrice)
	m := decimal.NewFromFloat(margin)

	spent := bp.Add(bp.Div(decimal.NewFromInt(100)).Mul(c).Round(2))
	gm := spent.Add(spent.Div(decimal.NewFromInt(100)).Mul(m).Round(2))

	target, _ := gm.Add(gm.Div(decimal.NewFromInt(100)).Mul(c).Round(2)).Float64()

	return target
}

func getItemsForSale(slots []domain.Slot, price float64) []domain.Slot {
	result := []domain.Slot{}
	p := decimal.NewFromFloat(price)

	for _, slot := range slots {
		if decimal.NewFromFloat(slot.TargetPrice).LessThanOrEqual(p) {
			result = append(result, slot)
		}
	}

	return result
}
