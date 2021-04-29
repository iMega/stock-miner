package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestDealings_AddTransaction(t *testing.T) {
	db, close, err := helpers.CreateDB(func(ctx context.Context, tx *sql.Tx) error {
		return dealingsTable(ctx, tx)
	})
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := Storage{db: db}

	tr := domain.Transaction{
		Slot: domain.Slot{
			Email: "email@example.com",
			ID:    "id",
			StockItem: domain.StockItem{
				Ticker: "ticker",
				FIGI:   "figi",
			},
			StartPrice:  1,
			ChangePrice: 2,
		},
		BuyOrderID: "54321",
		BuyAt:      time.Now(),
	}

	ctx := contexkey.WithEmail(context.Background(), tr.Slot.Email)
	if err := s.buyTransaction(ctx, tr); err != nil {
		t.Fatal(err)
	}

	trs, err := s.Dealings(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(trs) != 1 {
		t.Errorf("failed getting transactions not equals 1, got %d", len(trs))
	}

	assert.Equal(t, trs[0].BuyAt.Unix(), tr.BuyAt.Unix())

	tr.BuyAt = time.Now()
	trs[0].BuyAt = tr.BuyAt

	tr.SellAt = time.Now()
	trs[0].SellAt = tr.SellAt

	assert.Equal(t, trs[0], tr)
}

func TestDealings_ConfirmTransaction(t *testing.T) {
	db, close, err := helpers.CreateDB(func(ctx context.Context, tx *sql.Tx) error {
		return dealingsTable(ctx, tx)
	})
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := Storage{db: db}

	tr := domain.Transaction{
		Slot: domain.Slot{
			Email: "email@example.com",
			ID:    "id",
			StockItem: domain.StockItem{
				Ticker: "ticker",
				FIGI:   "figi",
			},
			StartPrice:  1.10,
			ChangePrice: 1.20,
			BuyingPrice: 3,
			TargetPrice: 4,
			Profit:      1,

			Qty:         2,
			AmountSpent: 6,
		},
		BuyOrderID: "54321",
		BuyAt:      time.Now(),
	}

	ctx := contexkey.WithEmail(context.Background(), tr.Slot.Email)
	if err := s.buyTransaction(ctx, tr); err != nil {
		t.Fatal(err)
	}

	if err := s.ConfirmBuyTransaction(ctx, tr); err != nil {
		t.Fatal(err)
	}

	trs, err := s.Dealings(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(trs) != 1 {
		t.Errorf("failed getting transactions not equals 1, got %d", len(trs))
	}

	assert.Equal(t, trs[0].BuyAt.Unix(), tr.BuyAt.Unix())

	tr.BuyAt = time.Now()
	trs[0].BuyAt = tr.BuyAt

	tr.SellAt = time.Now()
	trs[0].SellAt = tr.SellAt

	assert.Equal(t, trs[0], tr)
}

func TestDealings_SellTransaction(t *testing.T) {
	db, close, err := helpers.CreateDB(func(ctx context.Context, tx *sql.Tx) error {
		return dealingsTable(ctx, tx)
	})
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := Storage{db: db}

	tr := domain.Transaction{
		Slot: domain.Slot{
			Email: "email@example.com",
			ID:    "id",
			StockItem: domain.StockItem{
				Ticker: "ticker",
				FIGI:   "figi",
			},
			StartPrice:  1.10,
			ChangePrice: 1.20,
			BuyingPrice: 3,
			TargetPrice: 4,
			Profit:      1,

			Qty:         2,
			AmountSpent: 6,
		},
		SalePrice:    4,
		AmountIncome: 8,
		BuyOrderID:   "54321",
		SellOrderID:  "12345",
		BuyAt:        time.Now(),
		Duration:     0,
		SellAt:       time.Now(),
	}

	ctx := contexkey.WithEmail(context.Background(), tr.Slot.Email)
	if err := s.buyTransaction(ctx, tr); err != nil {
		t.Fatal(err)
	}

	if err := s.ConfirmBuyTransaction(ctx, tr); err != nil {
		t.Fatal(err)
	}

	if err := s.SellTransaction(ctx, tr); err != nil {
		t.Fatal(err)
	}

	trs, err := s.Dealings(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(trs) != 1 {
		t.Errorf("failed getting transactions not equals 1, got %d", len(trs))
	}

	assert.Equal(t, trs[0].BuyAt.Unix(), tr.BuyAt.Unix())
	assert.Equal(t, trs[0].SellAt.Unix(), tr.SellAt.Unix())

	tr.BuyAt = time.Now()
	trs[0].BuyAt = tr.BuyAt

	tr.SellAt = time.Now()
	trs[0].SellAt = tr.SellAt

	assert.Equal(t, trs[0], tr)
}
