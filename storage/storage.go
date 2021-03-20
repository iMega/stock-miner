package storage

import (
	"context"
	"database/sql"
	"os"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	_ "github.com/mattn/go-sqlite3"
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

func CreateDatabase(name string) error {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(name, os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return err
	}

	query := `CREATE TABLE price (
        symbol VARCHAR(64) NOT NULL,
        create_at DATETIME NOT NULL,
        price DECIMAL(10,5) NULL,
        CONSTRAINT price PRIMARY KEY (symbol, create_at)
    )`

	if _, err = db.Exec(query); err != nil {
		return err
	}

	return db.Close()
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
