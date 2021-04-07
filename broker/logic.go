package broker

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/robfig/cron/v3"
)

func (b *Broker) run() {
	inCh := make(chan domain.PriceReceiptMessage)
	outCh := make(chan domain.PriceReceiptMessage)
	psCh := make(chan domain.PriceReceiptMessage)
	w1 := b.makePricerChannel(inCh, outCh, psCh)
	w2 := b.makePriceStorageChannel(psCh)
	b.noName(outCh)

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
				res, err := b.Pricer.GetPrice(context.Background(), t)
				if err != nil {
					b.logger.Errorf("failed getting price from YF, %s", err)
				}
				res.Error = err

				// b.logger.Infof("-------- %s %f", res.Ticker, res.Price)

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
					b.logger.Infof("+++++222 %s", err)
				}
			})
		}
	}()

	return wp
}

func (b *Broker) noName(in chan domain.PriceReceiptMessage) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for task := range in {
			t := task
			wp.Submit(func() {
				ctx := contexkey.WithEmail(context.Background(), t.Email)
				settings, err := b.SettingsStorage.Settings(ctx)
				if err != nil {
					return
				}

				slots, err := b.Stack.Slot(ctx, t.FIGI)
				if err != nil {
					return
				}

				if settings.Slot.Volume <= len(slots) {
					return
				}

				ctx = contexkey.WithToken(ctx, settings.MarketCredentials[settings.MarketProvider].Token)
				ctx = contexkey.WithAPIURL(ctx, settings.MarketCredentials[settings.MarketProvider].APIURL)
				ob, _ := b.Market.OrderBook(ctx, t.StockItem)

				trend, err := b.SMAStack.IsTrendUp(t.Ticker)
				if err != nil {
					return
				}
				frame, err := b.SMAStack.Get(t.Ticker)
				if err != nil {
					return
				}

				// b.Stack.Add(t.Ticker, t.Price)
				// v, _ := b.Stack.Get(t.Ticker)
				b.logger.Infof("%s LP: %f, SMALP: %f, TrendUP: %v, FR:%#v", t.Ticker, ob.LastPrice, frame.Last, trend, frame.Fifo)
			})
		}
	}()

	return wp
}
