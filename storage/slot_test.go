package storage

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/tests/helpers"
)

func TestSlot_AddSlot(t *testing.T) {
	db, close, err := helpers.CreateDB(func(ctx context.Context, tx *sql.Tx) error {
		return slotTable(ctx, tx)
	})
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := Storage{
		db: db,
	}

	slot := domain.Slot{
		Email: "test@example.com",
		StockItem: domain.StockItem{
			FIGI: "figi",
		},
	}

	ctx := contexkey.WithEmail(context.Background(), slot.Email)
	if err := s.addSlot(ctx, slot); err != nil {
		t.Error(err)
	}

	slots, err := s.Slot(ctx, slot.FIGI)
	if err != nil {
		t.Error(err)
	}

	if len(slots) != 1 {
		t.Errorf("failed getting slots not equals 1, got %d", len(slots))
	}

	if !reflect.DeepEqual(slots[0], slot) {
		t.Errorf("storage.addSlot = %v, want %v", slots[0], slot)
	}
}
