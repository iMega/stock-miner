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
	b.makePricerChannel(inCh, outCh)

	delay := cron.DelayIfStillRunning(&logger{log: b.logger})
	b.cron.AddJob("@every 1s", delay(cron.FuncJob(func() {
		b.StockStorage.StockItemApprovedAll(context.Background(), inCh)
	})))
}

func (b *Broker) makePricerChannel(in, out chan domain.PriceReceiptMessage) {
	wp := workerpool.New(5)

	go func() {
		for t := range in {
			wp.Submit(func() {
				res, err := b.Pricer.GetPrice(context.Background(), t)
				res.Error = err
				out <- res
			})
		}
	}()
}
