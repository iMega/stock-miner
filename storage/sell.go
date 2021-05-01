package storage

import (
	"context"
	"database/sql"

	"github.com/imega/stock-miner/domain"
	tools "github.com/imega/stock-miner/sql"
	"github.com/imega/stock-miner/uuid"
)

func (s *Storage) Sell(ctx context.Context, t domain.Transaction) error {
	wrapper := tools.TxWrapper{s.db}
	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.SellTransaction(ctx, t); err != nil {
			return err
		}

		if err := s.deleteSlot(ctx, t.Slot); err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) ConfirmSell(ctx context.Context, t domain.Transaction) error {
	wrapper := tools.TxWrapper{s.db}
	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		return s.SellTransaction(ctx, t)
	})
}

func (s *Storage) PartialSell(ctx context.Context, t domain.Transaction, qty int) error {
	wrapper := tools.TxWrapper{s.db}
	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.SellTransaction(ctx, t); err != nil {
			return err
		}

		if err := s.deleteSlot(ctx, t.Slot); err != nil {
			return err
		}

		newTr := t

		newTr.Slot.ID = uuid.NewID().String()
		newTr.Slot.Qty = qty - t.Slot.Qty
		newTr.SalePrice = 0
		newTr.SellOrderID = ""
		newTr.TargetAmount = 0
		newTr.AmountIncome = 0

		if err := s.addSlot(ctx, newTr.Slot); err != nil {
			return err
		}

		if err := s.buyTransaction(ctx, newTr); err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) PartialConfirmSell(ctx context.Context, t domain.Transaction, qty int) error {
	wrapper := tools.TxWrapper{s.db}
	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		return s.SellTransaction(ctx, t)
	})
}
