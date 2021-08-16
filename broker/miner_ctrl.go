package broker

import (
	"context"
	"fmt"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
)

func (b *Broker) Stop() bool {
	b.cron.Stop()
	b.cronIsRunning = false

	if err := b.updateMainSettings(b.cronIsRunning); err != nil {
		b.logger.Errorf("miner is stopped, %s", err)
		b.cron.Stop()
	}

	b.logger.Infof("stop mining")

	return true
}

func (b *Broker) Start() bool {
	b.cron.Start()
	b.cronIsRunning = true

	if err := b.updateMainSettings(b.cronIsRunning); err != nil {
		b.logger.Errorf("miner is stopped, %s", err)
		b.cron.Stop()
	}

	b.logger.Infof("start mining")

	return true
}

func (b *Broker) updateMainSettings(status bool) error {
	ctx := contexkey.WithEmail(context.Background(), "main-settings")
	s := domain.Settings{
		MainSettings: domain.MainSettings{
			MiningStatus: status,
		},
	}

	if err := b.SettingsStorage.SaveSettings(ctx, s); err != nil {
		return fmt.Errorf("failed to update status global miner, %s", err)
	}

	return nil
}

func (b *Broker) Status() bool {
	return b.cronIsRunning
}

func (b *Broker) MainSettings() (domain.MainSettings, error) {
	ctx := contexkey.WithEmail(context.Background(), "main-settings")
	s, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		return domain.MainSettings{},
			fmt.Errorf("failed getting settings, %s", err)
	}

	return s.MainSettings, nil
}
