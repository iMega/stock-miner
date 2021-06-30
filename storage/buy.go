package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imega/stock-miner/domain"
	tools "github.com/imega/stock-miner/sql"
)

func (s *Storage) Buy(ctx context.Context, t domain.Transaction) error {
	tx := func(ctx context.Context, tx *sql.Tx) error {
		if err := s.addSlot(ctx, t.Slot); err != nil {
			return err
		}

		return s.buyTransaction(ctx, t)
	}

	wrapper := tools.TxWrapper{DB: s.db}
	if err := wrapper.Transaction(ctx, nil, tx); err != nil {
		return fmt.Errorf("failed to execute transaction, %w", err)
	}

	return nil
}

func (s *Storage) ConfirmBuy(ctx context.Context, t domain.Transaction) error {
	tx := func(ctx context.Context, tx *sql.Tx) error {
		if err := s.updateSlot(ctx, t.Slot); err != nil {
			return err
		}

		return s.ConfirmBuyTransaction(ctx, t)
	}

	wrapper := tools.TxWrapper{DB: s.db}
	if err := wrapper.Transaction(ctx, nil, tx); err != nil {
		return fmt.Errorf("failed to execute transaction, %w", err)
	}

	return nil
}
