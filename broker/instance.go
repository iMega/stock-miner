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
	Stack         smaStack
}

// New creates a new instance of Broker
func New(opts ...Option) *Broker {
	b := &Broker{
		cron:  cron.New(),
		Stack: make(smaStack),
	}

	for _, opt := range opts {
		opt(b)
	}

	b.run()

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
