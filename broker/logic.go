package broker

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/imega/stock-miner/domain"
	"github.com/robfig/cron/v3"
)

func (b *Broker) run() {
	inCh := make(chan domain.PriceReceiptMessage)
	outCh := make(chan domain.PriceReceiptMessage)
	psCh := make(chan domain.PriceReceiptMessage, 100)
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
		for t := range in {
			wp.Submit(func() {

				res, err := b.Pricer.GetPrice(context.Background(), t)
				if err != nil {
					b.logger.Errorf("failed getting price from YF, %s", err)
				}
				res.Error = err
				// b.logger.Info("+++++ %#v", t)
				// b.logger.Infof("++ %f", t.Price)

				if !b.Stack.Add(res.Ticker, res.Price) {
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
		for t := range in {
			b.logger.Info("+++++ %#v", t)
			wp.Submit(func() {
				err := b.StockStorage.AddMarketPrice(context.Background(), t)
				b.logger.Info("+++++222 %s", err)
			})
		}
	}()

	return wp
}

func (b *Broker) noName(in chan domain.PriceReceiptMessage) *workerpool.WorkerPool {
	wp := workerpool.New(1)

	go func() {
		for t := range in {
			wp.Submit(func() {
				// b.Stack.Add(t.Ticker, t.Price)
				v, _ := b.Stack.Get(t.Ticker)
				b.logger.Info("_____ %#v", v)
			})
		}
	}()

	return wp
}
