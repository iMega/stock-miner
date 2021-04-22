package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/tests/helpers"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, slots[0], slot)
}

func TestStorage_Slot(t *testing.T) {
	type args struct {
		ctx  context.Context
		figi string
	}

	db, close, err := helpers.CreateDB(func(ctx context.Context, tx *sql.Tx) error {
		return slotTable(ctx, tx)
	})
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	slot := domain.Slot{
		Email: "test@example.com",
		StockItem: domain.StockItem{
			FIGI: "figi",
		},
	}
	ctx := contexkey.WithEmail(context.Background(), slot.Email)

	s := Storage{db: db}

	if err := s.addSlot(ctx, slot); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		args    args
		want    []domain.Slot
		wantErr bool
	}{
		{
			name:    "with figi",
			args:    args{ctx: ctx, figi: "figi"},
			want:    []domain.Slot{slot},
			wantErr: false,
		},
		{
			name:    "without figi",
			args:    args{ctx: ctx, figi: ""},
			want:    []domain.Slot{slot},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Slot(tt.args.ctx, tt.args.figi)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Slot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got, tt.want)
		})
	}
}
