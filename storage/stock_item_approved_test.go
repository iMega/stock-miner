package storage

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	tools "github.com/imega/stock-miner/sql"
	"github.com/imega/stock-miner/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_stockItemApprovedTableFieldStartTime(t *testing.T) {
	db, close, err := helpers.CreateDB(stockItemApprovedCreateTableHelper)
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	expected, err := getTableInfo(db, "stock_item_approved")
	if err != nil {
		t.Fatalf("failed getting table info, %s", err)
	}

	expected.Columns = append(expected.Columns, col{
		CID:          6,
		Name:         "startTime",
		Type:         "INTEGER",
		NotNull:      1,
		DefaultValue: "11",
		PK:           0,
	})

	ctx := context.Background()
	wrapper := tools.TxWrapper{db}

	err = wrapper.Transaction(ctx, nil, stockItemApprovedTableFieldStartTime)
	if err != nil {
		t.Errorf("failed to migrate table, %s", err)
	}

	actual, err := getTableInfo(db, "stock_item_approved")
	if err != nil {
		t.Fatalf("failed getting table info, %s", err)
	}

	assert.Equal(t, expected, actual)
}

func stockItemApprovedCreateTableHelper(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS stock_item_approved (
        email VARCHAR(64) NOT NULL,
        ticker VARCHAR(64) NOT NULL,
        figi VARCHAR(200) NOT NULL,
        amount_limit FLOAT NOT NULL,
        transaction_limit INTEGER NOT NULL,
        currency VARCHAR(64) NOT NULL,
        CONSTRAINT pair PRIMARY KEY (email, ticker)
    )`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to create table stock_item_approved, %w", err)
	}

	return nil
}
