package storage

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	tools "github.com/imega/stock-miner/sql"
	_ "github.com/mattn/go-sqlite3"
)

type closeDB func() error

func createDB() (*sql.DB, closeDB, error) {
	file, err := ioutil.TempFile("", "stockminer")
	if err != nil {
		return nil, nil, err
	}

	filename := file.Name()
	if err := file.Close(); err != nil {
		return nil, nil, err
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	wrapper := tools.TxWrapper{db}
	wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		if err := slotTable(ctx, tx); err != nil {
			return err
		}

		return nil
	})

	return db,
		closeDB(func() error {
			errDB := db.Close()
			if err := os.Remove(filename); err != nil || errDB != nil {
				return fmt.Errorf("%s, %s", errDB, err)
			}

			return nil
		}),
		nil
}

func TestSlot_AddSlot(t *testing.T) {
	db, close, err := createDB()
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
