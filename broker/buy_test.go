package broker

import (
	"context"
	"testing"

	"github.com/imega/stock-miner/domain"
)

func TestBroker_confirmBuyJob(t *testing.T) {
	type fields struct {
		StockStorage    domain.StockStorage
		Market          domain.Market
		SettingsStorage domain.SettingsStorage
	}
	type args struct {
		tr domain.Transaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "optimistic",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Broker{
				StockStorage:    tt.fields.StockStorage,
				Market:          tt.fields.Market,
				SettingsStorage: tt.fields.SettingsStorage,
			}
			if err := b.confirmBuyJob(tt.args.tr); (err != nil) != tt.wantErr {
				t.Errorf("Broker.confirmBuyJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type stockStorageStub struct{}

func (stockStorageStub) AddStockItemApproved(_ context.Context, _ domain.StockItem) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) UpdateStockItemApproved(_ context.Context, _ domain.StockItem) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) RemoveStockItemApproved(_ context.Context, _ domain.StockItem) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) StockItemApprovedAll(_ context.Context, _ chan domain.Message) {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) StockItemApproved(_ context.Context) ([]domain.StockItem, error) {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) Slot(_ context.Context, _ string) ([]domain.Slot, error) {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) Dealings(_ context.Context) ([]domain.Transaction, error) {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) Transaction(_ context.Context, _ string) (domain.Transaction, error) {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) AddMarketPrice(_ context.Context, _ domain.PriceReceiptMessage) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) Buy(_ context.Context, _ domain.Transaction) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) ConfirmBuy(_ context.Context, _ domain.Transaction) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) Sell(_ context.Context, _ domain.Transaction) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) ConfirmSell(_ context.Context, _ domain.Transaction) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) PartialSell(_ context.Context, _ domain.Transaction, _ int) error {
	panic("not implemented") // TODO: Implement
}

func (stockStorageStub) PartialConfirmSell(_ context.Context, _ domain.Transaction, _ int) error {
	panic("not implemented") // TODO: Implement
}
