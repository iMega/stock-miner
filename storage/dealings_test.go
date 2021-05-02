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

func TestStorage_Dealings(t *testing.T) {
	db, close, err := helpers.CreateDB(func(ctx context.Context, tx *sql.Tx) error {
		return dealingsTable(ctx, tx)
	})
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := Storage{db: db}
	email := "email@example.com"
	now := time.Now()
	expected := []domain.Transaction{
		{
			Slot: domain.Slot{
				StockItem: domain.StockItem{
					Ticker: "ticker1",
					FIGI:   "figi1",
				},
				Email:       email,
				ID:          "1",
				StartPrice:  101,
				ChangePrice: 102,
				BuyingPrice: 103,
				TargetPrice: 104,
				Profit:      105,
				Qty:         106,
				AmountSpent: 108,
				TotalProfit: 109,
			},
			SalePrice:    110,
			AmountIncome: 111,
			BuyOrderID:   "112",
			SellOrderID:  "113",
			BuyAt:        now,
			Duration:     114,
			SellAt:       now,
		},
		{
			Slot: domain.Slot{
				StockItem: domain.StockItem{
					Ticker: "ticker2",
					FIGI:   "figi2",
				},
				Email:       email,
				ID:          "2",
				StartPrice:  201,
				ChangePrice: 202,
				BuyingPrice: 203,
				TargetPrice: 204,
				Profit:      205,
				Qty:         206,
				AmountSpent: 208,
			},
			BuyOrderID: "212",
			BuyAt:      now,
			SellAt:     now,
		},
	}
	ctx := contexkey.WithEmail(context.Background(), email)
	for _, tr := range expected {
		if err := s.buyTransaction(ctx, tr); err != nil {
			t.Fatal(err)
		}

		if err := s.ConfirmBuyTransaction(ctx, tr); err != nil {
			t.Fatal(err)
		}
	}

	if err := s.SellTransaction(ctx, expected[0]); err != nil {
		t.Fatal(err)
	}

	actual, err := s.Dealings(ctx)
	if err != nil {
		t.Fatalf("Storage.Dealings() error = %v", err)
	}

	for k := range actual {
		actual[k].BuyAt = now
		actual[k].SellAt = now
	}

	assert.Equal(t, expected, actual)
}
