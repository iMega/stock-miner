package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imega/stock-miner/broker"
	"github.com/imega/stock-miner/contexkey"
)

func (s *Storage) StockItemApproved(ctx context.Context) ([]broker.StockItem, error) {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user from context")
	}

	q := `select ticker, figi, amount_limit, transaction_limit from stock_item_approved where email = ?`
	rows, err := s.db.QueryContext(ctx, q, email)
	if err != nil {
		return nil, fmt.Errorf("failed getting approved stock items, %s", err)
	}

	var result []broker.StockItem
	for rows.Next() {
		var (
			ticker, figi     string
			amountLimit      float64
			transactionLimit int
		)
		if err := rows.Scan(&ticker, &figi, &amountLimit, &transactionLimit); err != nil {
			return nil, fmt.Errorf("failed to scan approved stock item, %s", err)
		}
		result = append(result, broker.StockItem{
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

func (s *Storage) AddStockItemApproved(ctx context.Context, item broker.StockItem) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	q := `insert into stock_item_approved (email, ticker, figi, amount_limit, transaction_limit) values (?,?,?,?,?)`
	_, err := s.db.ExecContext(ctx, q, email, item.Ticker, item.FIGI, item.AmountLimit, item.TransactionLimit)
	if err != nil {
		return fmt.Errorf("failed to add approved stock item, %s", err)
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
        CONSTRAINT pair PRIMARY KEY (email, ticker)
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}
