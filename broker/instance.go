package broker

import (
	"github.com/imega/daemon"
	"github.com/imega/stock-miner/domain"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Broker is the main struct.
type Broker struct {
	StockStorage    domain.StockStorage
	Pricer          domain.Pricer
	Market          domain.Market
	SMAStack        domain.SMAStack
	SettingsStorage domain.SettingsStorage
	Stack           domain.Stack

	logger        logrus.FieldLogger
	cron          *cron.Cron
	isShutdown    bool
	cronIsRunning bool
}

// New creates a new instance of Broker.
func New(opts ...Option) *Broker {
	b := &Broker{
		cron:     cron.New(),
		SMAStack: NewSMAStack(),
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

func WithMarket(d domain.Market) Option {
	return func(b *Broker) {
		b.Market = d
	}
}

func WithSettingsStorage(d domain.SettingsStorage) Option {
	return func(b *Broker) {
		b.SettingsStorage = d
	}
}

func WithStack(d domain.Stack) Option {
	return func(b *Broker) {
		b.Stack = d
	}
}

func WithSMAStack(d domain.SMAStack) Option {
	return func(b *Broker) {
		b.SMAStack = d
	}
}
