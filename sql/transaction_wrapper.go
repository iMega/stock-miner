package sql

import (
	"context"
	"database/sql"
	"fmt"
)

type txWrapper struct{ *sql.DB }

type TxFunc func(context.Context, *sql.Tx) error

func (w *txWrapper) Transaction(
	ctx context.Context,
	opts *sql.TxOptions,
	f TxFunc,
) error {
	tx, err := w.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	if err := f(ctx, tx); err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("%s, failed to rollback transaction, %s", err, e)
		}

		return err
	}

	return tx.Commit()
}
