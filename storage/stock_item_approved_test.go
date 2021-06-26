package storage

import (
	"context"
	"database/sql"
	"testing"

	tools "github.com/imega/stock-miner/sql"
	"github.com/imega/stock-miner/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_stockItemApprovedTableFieldStartTime(t *testing.T) {
	db, close, err := helpers.CreateDB(
		func(ctx context.Context, tx *sql.Tx) error {
			return stockItemApprovedTable(ctx, tx)
		},
	)
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
