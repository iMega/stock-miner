package storage

import (
	"context"
	"database/sql"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

type Option func(b *Storage)

type Storage struct {
	db *sql.DB
}

func WithSqllite(db *sql.DB) Option {
	return func(s *Storage) {
		s.db = db
	}
}

func New(opts ...Option) *Storage {
	s := &Storage{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Storage) AddMarketPrice(ctx context.Context, o sdk.RestOrderBook) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var (
		price sql.NullFloat64
		t     sql.NullTime
	)

	getLastPrice := "select price, max(create_at) from price where symbol = ?"
	if err := tx.QueryRowContext(ctx, getLastPrice, o.FIGI).Scan(&price, &t); err != nil {
		return err
	}

	if price.Valid && price.Float64 == o.LastPrice {
		return nil
	}

	insertPrice := "insert into price (symbol, create_at, price) values (?, ?, ?)"
	_, err = tx.ExecContext(ctx, insertPrice, o.FIGI, time.Now().String(), o.LastPrice)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}
