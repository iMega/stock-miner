package broker

import (
	"context"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/daemon"
	"github.com/imega/stock-miner/yahooprovider"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Broker is the main struct
type Broker struct {
	storage    Storage
	logger     logrus.FieldLogger
	isShutdown bool
}

type Storage interface {
	AddMarketPrice(context.Context, sdk.RestOrderBook) error
}

// New creates a new instance of Broker
func New(opts ...Option) *Broker {
	b := &Broker{}

	for _, opt := range opts {
		opt(b)
	}

	l := &logger{log: b.logger}

	c := cron.New()

	delay := cron.DelayIfStillRunning(l)

	yProvider := yahooprovider.New()

	c.AddJob("@every 1s", delay(cron.FuncJob(func() {
		p, err := yProvider.Price(context.Background(), sdk.Instrument{
			FIGI: "AAPL",
		})
		if err != nil {
			b.logger.Errorf("failed task, %s", err)
		}

		if err := b.storage.AddMarketPrice(context.Background(), p); err != nil {
			b.logger.Errorf("failed to add price, %s", err)
		}

		b.logger.Infof("======= %#v", p)
	})))

	// c.Start()

	return b
}

type Option func(b *Broker)

func WithStorage(s Storage) Option {
	return func(b *Broker) {
		b.storage = s
	}
}

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
