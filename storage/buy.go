package storage

import (
	"context"
	"database/sql"

	"github.com/imega/stock-miner/domain"
	tools "github.com/imega/stock-miner/sql"
)

func (s *Storage) Buy(ctx context.Context, t domain.Transaction) error {
	wrapper := tools.TxWrapper{s.db}
	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.addSlot(ctx, t.Slot); err != nil {
			return err
		}

		if err := s.buyTransaction(ctx, t); err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) ConfirmBuy(ctx context.Context, t domain.Transaction) error {
	wrapper := tools.TxWrapper{s.db}
	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.updateSlot(ctx, t.Slot); err != nil {
			return err
		}

		if err := s.ConfirmBuyTransaction(ctx, t); err != nil {
			return err
		}

		return nil
	})
}
