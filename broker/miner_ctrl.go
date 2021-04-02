package broker

func (b *Broker) Stop() bool {
	b.cron.Stop()
	b.cronIsRunning = false

	return true
}

func (b *Broker) Start() bool {
	b.cron.Start()
	b.cronIsRunning = true

	return true
}

func (b *Broker) Status() bool {
	return b.cronIsRunning
}
