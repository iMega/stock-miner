package broker

import (
	"github.com/imega/daemon"
	"github.com/imega/stock-miner/domain"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Broker is the main struct
type Broker struct {
	StockStorage  domain.StockStorage
	Pricer        domain.Pricer
	logger        logrus.FieldLogger
	isShutdown    bool
	cron          *cron.Cron
	cronIsRunning bool
}

// New creates a new instance of Broker
func New(opts ...Option) *Broker {
	b := &Broker{
		cron: cron.New(),
	}

	for _, opt := range opts {
		opt(b)
	}

	// wp := workerpool.New(5)
	// task := make(chan domain.PriceReceiptMessageOut)

	// go func() error {
	// 	for {
	// 		select {
	// 		case t, ok := <-task:
	// 			if !ok {
	// 				return nil
	// 			}

	// 			wp.Submit(func() {
	// 				// fmt.Printf("--START----%s-%s\n", t.Ticker, t.MarketState)
	// 				<-time.After(5 * time.Second)
	// 				fmt.Printf("--END------%s-%s\n", t.Ticker, t.MarketState)
	// 			})
	// 		}
	// 	}
	// }()

	b.run()

	// b.cron.AddJob("@every 1s", delay(cron.FuncJob(func() {

	// 	i++

	// 	if 20 == wp.WaitingQueueSize() {
	// 		c.Stop()
	// 		<-time.After(10 * time.Second)
	// 		c.Start()
	// 	}

	// 	fmt.Printf("-------------------%d----------%d\n", i, wp.WaitingQueueSize())

	// 	b.StockStorage.StockItemApprovedAll(context.Background(), task)
	// 	// p, err := yProvider.Price(context.Background(), sdk.Instrument{
	// 	// 	FIGI: "AAPL",
	// 	// })
	// 	// if err != nil {
	// 	// 	b.logger.Errorf("failed task, %s", err)
	// 	// }

	// 	// if err := b.storage.AddMarketPrice(context.Background(), p); err != nil {
	// 	// 	b.logger.Errorf("failed to add price, %s", err)
	// 	// }
	// 	// task <- helpers.RandomInt(1, 70)
	// 	// b.logger.Infof("======= %#v", p)
	// })))

	return b
}

type Option func(b *Broker)

func WithLogger(l logrus.FieldLogger) Option {
	return func(b *Broker) {
		b.logger = l
	}
}

func (b *Broker) ShutdownFunc() daemon.ShutdownFunc {
	return func() {
		b.isShutdown = true
	}
}

func WithStockStorage(s domain.StockStorage) Option {
	return func(b *Broker) {
		b.StockStorage = s
	}
}

func WithPricer(p domain.Pricer) Option {
	return func(b *Broker) {
		b.Pricer = p
	}
}
