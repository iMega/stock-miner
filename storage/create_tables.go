package storage

import (
	"context"
	"database/sql"
	"os"

	tools "github.com/imega/stock-miner/sql"
)

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

	ctx := context.Background()
	wrapper := tools.TxWrapper{db}
	wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		if err := priceTable(ctx, tx); err != nil {
			return err
		}

		if err := userTable(ctx, tx); err != nil {
			return err
		}

		if err := stockItemApprovedTable(ctx, tx); err != nil {
			return err
		}

		if err := settingsTable(ctx, tx); err != nil {
			return err
		}

		if err := dealingsTable(ctx, tx); err != nil {
			return err
		}

		if err := slotTable(ctx, tx); err != nil {
			return err
		}

		return nil
	})

	return db.Close()
}

func priceTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS price (
        symbol VARCHAR(64) NOT NULL,
        create_at DATETIME NOT NULL,
        price DECIMAL(10,5) NULL,
        CONSTRAINT price PRIMARY KEY (symbol, create_at)
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}

func userTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS user (
        email VARCHAR(64) PRIMARY KEY,
        name VARCHAR(64),
        avatar VARCHAR(200),
        id VARCHAR(200),
        deleted INTEGER DEFAULT 0,
        role CHAR(4),
        create_at DATETIME NOT NULL
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}
