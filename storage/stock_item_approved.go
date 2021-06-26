package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
)

func (s *Storage) StockItemApproved(ctx context.Context) ([]domain.StockItem, error) {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user from context")
	}

	q := `select ticker, figi, amount_limit, transaction_limit from stock_item_approved where email = ?`
	rows, err := s.db.QueryContext(ctx, q, email)
	if err != nil {
		return nil, fmt.Errorf("failed getting approved stock items, %s", err)
	}

	var result []domain.StockItem
	for rows.Next() {
		var (
			ticker, figi     string
			amountLimit      float64
			transactionLimit int
		)
		if err := rows.Scan(&ticker, &figi, &amountLimit, &transactionLimit); err != nil {
			return nil, fmt.Errorf("failed to scan approved stock item, %s", err)
		}
		result = append(result, domain.StockItem{
			Ticker:           ticker,
			FIGI:             figi,
			AmountLimit:      amountLimit,
			TransactionLimit: transactionLimit,
		})
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows approved stock item, %s", err)
	}

	return result, nil
}

func (s *Storage) StockItemApprovedAll(
	ctx context.Context,
	out chan domain.Message,
) {
	query := `select email, ticker, figi, amount_limit, transaction_limit, currency from stock_item_approved`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		out <- domain.Message{
			Error: fmt.Errorf("failed getting approved stock items, %s", err),
		}

		return
	}

	for rows.Next() {
		var (
			email, ticker, figi string
			currency            string
			amountLimit         float64
			transactionLimit    int
		)

		if err := rows.Scan(&email, &ticker, &figi, &amountLimit, &transactionLimit, &currency); err != nil {
			out <- domain.Message{
				Error: fmt.Errorf("failed to scan approved stock item, %s", err),
			}

			return
		}

		out <- domain.Message{
			Transaction: domain.Transaction{
				Slot: domain.Slot{
					Email: email,
					StockItem: domain.StockItem{
						Ticker:           ticker,
						FIGI:             figi,
						AmountLimit:      amountLimit,
						TransactionLimit: transactionLimit,
						Currency:         currency,
					},
				},
			},
		}
	}

	if err := rows.Close(); err != nil {
		out <- domain.Message{
			Error: fmt.Errorf("failed to close rows approved stock item, %s", err),
		}
	}
}

func (s *Storage) AddStockItemApproved(ctx context.Context, item domain.StockItem) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	q := `insert into stock_item_approved (email, ticker, figi, amount_limit, transaction_limit, currency) values (?,?,?,?,?,?)`
	_, err := s.db.ExecContext(ctx, q, email, item.Ticker, item.FIGI, item.AmountLimit, item.TransactionLimit, item.Currency)
	if err != nil {
		return fmt.Errorf("failed to add approved stock item, %s", err)
	}

	return nil
}

func (s *Storage) UpdateStockItemApproved(ctx context.Context, item domain.StockItem) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	q := `update stock_item_approved set amount_limit=?, transaction_limit=? where email=? and ticker=?`
	_, err := s.db.ExecContext(ctx, q, item.AmountLimit, item.TransactionLimit, email, item.Ticker)
	if err != nil {
		return fmt.Errorf("failed to update approved stock item, %s", err)
	}

	return nil
}

func (s *Storage) RemoveStockItemApproved(ctx context.Context, item domain.StockItem) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	q := `delete from stock_item_approved where email=? and ticker=?`
	_, err := s.db.ExecContext(ctx, q, email, item.Ticker)
	if err != nil {
		return fmt.Errorf("failed to delete approved stock item, %s", err)
	}

	return nil
}

func stockItemApprovedTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS stock_item_approved (
        email VARCHAR(64) NOT NULL,
        ticker VARCHAR(64) NOT NULL,
        figi VARCHAR(200) NOT NULL,
        amount_limit FLOAT NOT NULL,
        transaction_limit INTEGER NOT NULL,
        currency VARCHAR(64) NOT NULL,
        CONSTRAINT pair PRIMARY KEY (email, ticker)
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}

func stockItemApprovedTableMigrate(ctx context.Context, tx *sql.Tx, ti tableInfo) error {
	if !hasColumn(ti, col{Name: "startTime"}) {
		if err := stockItemApprovedTableFieldStartTime(ctx, tx); err != nil {
			return fmt.Errorf("failed to migrate table stock_item_approved, %s", err)
		}
	}

	if !hasColumn(ti, col{Name: "endTime"}) {
		if err := stockItemApprovedTableFieldEndTime(ctx, tx); err != nil {
			return fmt.Errorf("failed to migrate table stock_item_approved, %s", err)
		}
	}

	return nil
}

func stockItemApprovedTableFieldStartTime(ctx context.Context, tx *sql.Tx) error {
	q := `ALTER TABLE stock_item_approved ADD startTime INTEGER NOT NULL DEFAULT 11`

	_, err := tx.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf(
			"failed to execute stockItemApprovedTableFieldStartTime, %s",
			err,
		)
	}

	return nil
}

func stockItemApprovedTableFieldEndTime(ctx context.Context, tx *sql.Tx) error {
	q := `ALTER TABLE stock_item_approved ADD endTime INTEGER NOT NULL DEFAULT 20`

	_, err := tx.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf(
			"failed to execute stockItemApprovedTableFieldEndTime, %s",
			err,
		)
	}

	return nil
}
