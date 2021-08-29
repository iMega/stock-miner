package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
)

func (s *Storage) StockItemApproved(
	ctx context.Context,
) ([]domain.StockItem, error) {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return nil, contexkey.ErrExtractEmail
	}

	q := `select ticker,
            figi,
            currency,
            amount_limit,
            transaction_limit,
            startTime,
            endTime,
            active,
            max_price
        from stock_item_approved
        where email = ?`

	rows, err := s.db.QueryContext(ctx, q, email)
	if err != nil {
		return nil, fmt.Errorf("failed getting approved stock items, %w", err)
	}
	defer rows.Close()

	var result []domain.StockItem

	for rows.Next() {
		var item domain.StockItem

		err := rows.Scan(
			&item.Ticker,
			&item.FIGI,
			&item.Currency,
			&item.AmountLimit,
			&item.TransactionLimit,
			&item.StartTime,
			&item.EndTime,
			&item.IsActive,
			&item.MaxPrice,
		)
		if err != nil {
			return nil,
				fmt.Errorf("failed to scan approved stock item, %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed getting row, %w", err)
	}

	return result, nil
}

func (s *Storage) StockItemApprovedAll(
	ctx context.Context,
	out chan domain.Message,
) {
	query := `select email,
                    ticker,
                    figi,
                    amount_limit,
                    transaction_limit,
                    currency,
                    startTime,
                    endTime,
                    max_price
                from stock_item_approved
                where active=1`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		out <- domain.Message{
			Error: fmt.Errorf(
				"failed getting all approved stock items, %s",
				err,
			),
		}

		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			email string
			item  domain.StockItem
		)

		err := rows.Scan(
			&email,
			&item.Ticker,
			&item.FIGI,
			&item.AmountLimit,
			&item.TransactionLimit,
			&item.Currency,
			&item.StartTime,
			&item.EndTime,
			&item.MaxPrice,
		)
		if err != nil {
			out <- domain.Message{
				Error: fmt.Errorf(
					"failed to scan approved stock item, %s",
					err,
				),
			}

			return
		}

		if !isValidPeriod(int(item.StartTime), int(item.EndTime)) {
			continue
		}

		out <- domain.Message{
			Transaction: domain.Transaction{
				Slot: domain.Slot{
					Email:     email,
					StockItem: item,
				},
			},
		}
	}

	if err := rows.Err(); err != nil {
		out <- domain.Message{
			Error: fmt.Errorf("failed getting row, %w", err),
		}

		return
	}
}

func isValidPeriod(start, end int) bool {
	currentHour := time.Now().Hour()

	return start <= currentHour && currentHour <= end
}

func (s *Storage) AddStockItemApproved(
	ctx context.Context,
	item domain.StockItem,
) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return contexkey.ErrExtractEmail
	}

	q := `insert into stock_item_approved (
            email,
            ticker,
            figi,
            amount_limit,
            transaction_limit,
            currency,
            startTime,
            endTime,
            active,
            max_price
        )
        values (?,?,?,?,?,?,?,?,?,?)`

	_, err := s.db.ExecContext(
		ctx,
		q,
		email,
		item.Ticker,
		item.FIGI,
		item.AmountLimit,
		item.TransactionLimit,
		item.Currency,
		item.StartTime,
		item.EndTime,
		item.IsActive,
		item.MaxPrice,
	)
	if err != nil {
		return fmt.Errorf("failed to add approved stock item, %w", err)
	}

	return nil
}

func (s *Storage) UpdateStockItemApproved(
	ctx context.Context,
	item domain.StockItem,
) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return contexkey.ErrExtractEmail
	}

	q := `update stock_item_approved
        set amount_limit=?,
            transaction_limit=?,
            startTime=?,
            endTime=?,
            active=?,
            max_price=?
        where email=? and ticker=?`

	_, err := s.db.ExecContext(
		ctx,
		q,
		item.AmountLimit,
		item.TransactionLimit,
		item.StartTime,
		item.EndTime,
		item.IsActive,
		item.MaxPrice,
		email,
		item.Ticker,
	)
	if err != nil {
		return fmt.Errorf("failed to update approved stock item, %w", err)
	}

	return nil
}

func (s *Storage) UpdateActiveStatusStockItemApproved(
	ctx context.Context,
	status bool,
) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return contexkey.ErrExtractEmail
	}

	q := `update stock_item_approved
        set active=?
        where email=?`

	if _, err := s.db.ExecContext(ctx, q, status, email); err != nil {
		return fmt.Errorf("failed to update approved stock item, %w", err)
	}

	return nil
}

func (s *Storage) RemoveStockItemApproved(
	ctx context.Context,
	item domain.StockItem,
) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return contexkey.ErrExtractEmail
	}

	q := `delete from stock_item_approved where email=? and ticker=?`
	if _, err := s.db.ExecContext(ctx, q, email, item.Ticker); err != nil {
		return fmt.Errorf("failed to delete approved stock item, %w", err)
	}

	return nil
}

func stockItemApprovedCreateTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS stock_item_approved (
        email VARCHAR(64) NOT NULL,
        ticker VARCHAR(64) NOT NULL,
        figi VARCHAR(200) NOT NULL,
        amount_limit FLOAT NOT NULL,
        transaction_limit INTEGER NOT NULL,
        currency VARCHAR(64) NOT NULL,
        startTime INTEGER NOT NULL DEFAULT 11,
        endTime INTEGER NOT NULL DEFAULT 20,
        active INTEGER NOT NULL DEFAULT 1,
        max_price INTEGER NOT NULL DEFAULT 0,
        CONSTRAINT pair PRIMARY KEY (email, ticker)
    )`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to create table stock_item_approved, %w", err)
	}

	return nil
}

func stockItemApprovedTableMigrate(
	ctx context.Context,
	tx *sql.Tx,
	ti tableInfo,
) error {
	if !hasColumn(ti, col{Name: "startTime"}) {
		if err := stockItemApprovedTableFieldStartTime(ctx, tx); err != nil {
			return fmt.Errorf(
				"failed to migrate table stock_item_approved, %w",
				err,
			)
		}
	}

	if !hasColumn(ti, col{Name: "endTime"}) {
		if err := stockItemApprovedTableFieldEndTime(ctx, tx); err != nil {
			return fmt.Errorf(
				"failed to migrate table stock_item_approved, %w",
				err,
			)
		}
	}

	if !hasColumn(ti, col{Name: "active"}) {
		if err := stockItemApprovedTableFieldActive(ctx, tx); err != nil {
			return fmt.Errorf(
				"failed to migrate table stock_item_approved, %w",
				err,
			)
		}
	}

	if !hasColumn(ti, col{Name: "max_price"}) {
		if err := stockItemApprovedTableFieldMaxPrice(ctx, tx); err != nil {
			return fmt.Errorf(
				"failed to migrate table stock_item_approved, %w",
				err,
			)
		}
	}

	return nil
}

func stockItemApprovedTableFieldStartTime(
	ctx context.Context,
	tx *sql.Tx,
) error {
	q := `ALTER TABLE stock_item_approved
        ADD startTime INTEGER NOT NULL DEFAULT 11`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf(
			"failed to execute stockItemApprovedTableFieldStartTime, %w",
			err,
		)
	}

	return nil
}

func stockItemApprovedTableFieldEndTime(
	ctx context.Context,
	tx *sql.Tx,
) error {
	q := `ALTER TABLE stock_item_approved
        ADD endTime INTEGER NOT NULL DEFAULT 20`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf(
			"failed to execute stockItemApprovedTableFieldEndTime, %w",
			err,
		)
	}

	return nil
}

func stockItemApprovedTableFieldActive(ctx context.Context, tx *sql.Tx) error {
	q := `ALTER TABLE stock_item_approved ADD active INTEGER NOT NULL DEFAULT 1`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf(
			"failed to execute stockItemApprovedTableFieldActive, %w",
			err,
		)
	}

	return nil
}

func stockItemApprovedTableFieldMaxPrice(ctx context.Context, tx *sql.Tx) error {
	q := `ALTER TABLE stock_item_approved ADD max_price INTEGER NOT NULL DEFAULT 0`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf(
			"failed to execute stockItemApprovedTableFieldMaxPrice, %w",
			err,
		)
	}

	return nil
}
