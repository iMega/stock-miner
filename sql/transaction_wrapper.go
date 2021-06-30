package sql

import (
	"context"
	"database/sql"
	"fmt"
)

type TxWrapper struct{ DB *sql.DB }

type TxFunc func(context.Context, *sql.Tx) error

func (w *TxWrapper) Transaction(
	ctx context.Context,
	opts *sql.TxOptions,
	f TxFunc,
) error {
	tx, err := w.DB.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction, %w", err)
	}

	if err := f(ctx, tx); err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("failed to execute transaction, %w", err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction, %w", err)
	}

	return nil
}
