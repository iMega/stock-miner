package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/imega/stock-miner/domain"
)

type Option func(b *Storage)

type Storage struct {
	db       *sql.DB
	settings map[string]domain.Settings
}

func WithSqllite(db *sql.DB) Option {
	return func(s *Storage) {
		s.db = db
	}
}

func New(opts ...Option) *Storage {
	s := &Storage{
		settings: make(map[string]domain.Settings),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Storage) AddMarketPrice(ctx context.Context, msg domain.PriceReceiptMessage) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction, %w", err)
	}

	// var (
	// 	price sql.NullFloat64
	// 	t     sql.NullString
	// )

	// getLastPrice := "select price, max(create_at) as last_time from price where symbol = ?"
	// if err := tx.QueryRowContext(ctx, getLastPrice, msg.Ticker).Scan(&price, &t); err != nil {
	// 	fmt.Printf("==ERRR= %s\n", err)
	// 	if err != sql.ErrNoRows {
	// 		return err
	// 	}
	// 	price = sql.NullFloat64{Valid: true}
	// }

	// if price.Valid && price.Float64 == msg.Price {
	// 	return nil
	// }

	q := "insert into price (symbol, create_at, price) values (?, ?, ?)"

	_, err = tx.ExecContext(ctx, q, msg.Ticker, time.Now().String(), msg.Price)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			return fmt.Errorf("failed to execute query, %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction, %w", err)
	}

	return nil
}
