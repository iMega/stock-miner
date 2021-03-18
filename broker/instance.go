package broker

import "github.com/imega/daemon"

// Broker is the main struct
type Broker struct {
	storage    Storage
	isShutdown bool
}

type Storage interface{}

// New creates a new instance of Broker
func New(opts ...Option) *Broker {
	b := &Broker{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

type Option func(b *Broker)

func WithStorage(s Storage) Option {
	return func(b *Broker) {
		b.storage = s
	}
}

func (b *Broker) ShutdownFunc() daemon.ShutdownFunc {
	return func() {
		b.isShutdown = true
	}
}
