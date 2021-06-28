package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	tools "github.com/imega/stock-miner/sql"
)

func CreateDatabase(name string) error {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(name, os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to create file database, %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close file database, %w", err)
	}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return fmt.Errorf("failed to open database, %w", err)
	}

	ctx := context.Background()
	wrapper := tools.TxWrapper{db}
	tx := func(ctx context.Context, tx *sql.Tx) error {
		if err := priceTable(ctx, tx); err != nil {
			return err
		}

		if err := userTable(ctx, tx); err != nil {
			return err
		}

		if err := stockItemApprovedCreateTable(ctx, tx); err != nil {
			return err
		}

		if err := settingsTable(ctx, tx); err != nil {
			return err
		}

		if err := dealingsTable(ctx, tx); err != nil {
			return err
		}

		return slotTable(ctx, tx)
	}

	if err := wrapper.Transaction(ctx, nil, tx); err != nil {
		return fmt.Errorf("failed to execute transaction, %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database, %w", err)
	}

	return nil
}

func priceTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS price (
        symbol VARCHAR(64) NOT NULL,
        create_at DATETIME NOT NULL,
        price DECIMAL(10,5) NULL,
        CONSTRAINT price PRIMARY KEY (symbol, create_at)
    )`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to execute query, %w", err)
	}

	return nil
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

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to execute query, %w", err)
	}

	return nil
}
