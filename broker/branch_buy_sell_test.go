package broker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/imega/stock-miner/domain"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBroker_branchBuyOrSell(t *testing.T) {
	type fields struct {
		StockStorage    domain.StockStorage
		Pricer          domain.Pricer
		Market          domain.Market
		SMAStack        domain.SMAStack
		SettingsStorage domain.SettingsStorage
		Stack           domain.Stack
		Traffic         StockItemTraffic
		logger          logrus.FieldLogger
		cron            *cron.Cron
		isShutdown      bool
		cronIsRunning   bool
		isDevMode       bool
	}
	type args struct {
		msg domain.Message
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		wantSell []domain.Slot
		wantBuy  []domain.Message
	}{
		// {
		// 	name: "if SMAframe does'nt exists then fail",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{IsGetStackFail: true}),
		// 	},
		// 	args:    args{},
		// 	wantErr: true,
		// },
		// {
		// 	name: "if range of SMAframe is zero then fail",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{SMAFrame: fakeSMAFrame{}}),
		// 	},
		// 	args:    args{},
		// 	wantErr: true,
		// },
		// {
		// 	name: "if SMAframe high range is zero then fail",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{
		// 			SMAFrame: fakeSMAFrame{RangeLow: 1},
		// 		}),
		// 	},
		// 	args:    args{},
		// 	wantErr: true,
		// },
		// {
		// 	name: "if SMAframe low range is zero then fail",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{
		// 			SMAFrame: fakeSMAFrame{RangeHigh: 1},
		// 		}),
		// 	},
		// 	args:    args{},
		// 	wantErr: true,
		// },
		// {
		// 	name: "if SMAframe is empty then fail",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{
		// 			SMAFrame: fakeSMAFrame{
		// 				RangeHigh: 1,
		// 				RangeLow:  1,
		// 				IsEmpty:   true,
		// 			},
		// 		}),
		// 	},
		// 	args:    args{},
		// 	wantErr: true,
		// },
		// {
		// 	name: "settings storage return fail",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{
		// 			SMAFrame: fakeSMAFrame{
		// 				RangeHigh: 1,
		// 				RangeLow:  1,
		// 				IsEmpty:   false,
		// 			},
		// 		}),
		// 		SettingsStorage: &fakeSettingsStorage{
		// 			IsFail: true,
		// 		},
		// 		logger: &fakeLogger{},
		// 	},
		// 	args: args{
		// 		msg: domain.Message{
		// 			Transaction: domain.Transaction{
		// 				Slot: domain.Slot{
		// 					Email: "info@example.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "slot storage return fail",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{
		// 			SMAFrame: fakeSMAFrame{
		// 				RangeHigh: 1,
		// 				RangeLow:  1,
		// 				IsEmpty:   false,
		// 			},
		// 		}),
		// 		SettingsStorage: &fakeSettingsStorage{
		// 			FakeSettings: domain.Settings{
		// 				//
		// 			},
		// 		},
		// 		Stack: &fakeStack{
		// 			GetSlotIsFail: true,
		// 		},
		// 		logger: &fakeLogger{},
		// 	},
		// 	args: args{
		// 		msg: domain.Message{
		// 			Transaction: domain.Transaction{
		// 				Slot: domain.Slot{
		// 					Email: "info@example.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "slot is empty",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{
		// 			SMAFrame: fakeSMAFrame{
		// 				RangeHigh: 1,
		// 				RangeLow:  1,
		// 				IsEmpty:   false,
		// 			},
		// 		}),
		// 		SettingsStorage: &fakeSettingsStorage{
		// 			FakeSettings: domain.Settings{},
		// 		},
		// 		Stack:  &fakeStack{},
		// 		logger: &fakeLogger{},
		// 	},
		// 	args: args{
		// 		msg: domain.Message{
		// 			Transaction: domain.Transaction{
		// 				Slot: domain.Slot{
		// 					Email: "info@example.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "slot will be sold",
		// 	fields: fields{
		// 		SMAStack: getSMAStack(fakeSMAStack{
		// 			SMAFrame: fakeSMAFrame{
		// 				RangeHigh: 1,
		// 				RangeLow:  1,
		// 				IsEmpty:   false,
		// 				LastValue: 10,
		// 				PrevValue: 10,
		// 			},
		// 		}),
		// 		SettingsStorage: &fakeSettingsStorage{
		// 			FakeSettings: domain.Settings{
		// 				//
		// 			},
		// 		},
		// 		Stack: &fakeStack{
		// 			Slots: []domain.Slot{
		// 				{
		// 					Email:       "info@example.com",
		// 					BuyingPrice: 10,
		// 				},
		// 			},
		// 		},
		// 		logger:  &fakeLogger{},
		// 		Traffic: NewStockItemTraffic(),
		// 	},
		// 	args: args{
		// 		msg: domain.Message{
		// 			Transaction: domain.Transaction{
		// 				Slot: domain.Slot{
		// 					Email: "info@example.com",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// 	wantSell: []domain.Slot{
		// 		{
		// 			Email:       "info@example.com",
		// 			BuyingPrice: 10,
		// 		},
		// 	},
		// },
		{
			name: "buy slot with max price option",
			fields: fields{
				SMAStack: getSMAStack(fakeSMAStack{
					SMAFrame: fakeSMAFrame{
						RangeHigh: 150,
						RangeLow:  50,
						IsEmpty:   false,
					},
				}),
				SettingsStorage: &fakeSettingsStorage{
					FakeSettings: domain.Settings{
						Slot: domain.SlotSettings{
							Volume: 1,
						},
					},
				},
				Stack: &fakeStack{
					Slots: []domain.Slot{},
				},
				logger:       &fakeLogger{},
				Market:       &fakeMarket{},
				StockStorage: &fakeStockStorage{},
				Traffic:      NewStockItemTraffic(),
			},
			args: args{
				msg: domain.Message{
					Transaction: domain.Transaction{
						Slot: domain.Slot{
							Email: "info@example.com",
							StockItem: domain.StockItem{
								MaxPrice: 100,
							},
						},
					},
					Price: 101,
				},
			},
			wantErr: false,
			wantBuy: []domain.Message{
				{
					Transaction: domain.Transaction{
						Slot: domain.Slot{
							Email: "info@example.com",
							StockItem: domain.StockItem{
								MaxPrice: 100,
							},
							SlotID: 1,
							Qty:    1,
						},
					},
				},
			},
		},
		{
			name: "if the maxPrice is equal to or less than the current price then the app shouldn’t complete the purchase",
			fields: fields{
				SMAStack: getSMAStack(fakeSMAStack{
					SMAFrame: fakeSMAFrame{
						RangeHigh: 150,
						RangeLow:  50,
						IsEmpty:   false,
					},
				}),
				SettingsStorage: &fakeSettingsStorage{
					FakeSettings: domain.Settings{
						Slot: domain.SlotSettings{
							Volume: 1,
						},
					},
				},
				Stack: &fakeStack{
					Slots: []domain.Slot{},
				},
				logger:       &fakeLogger{},
				Market:       &fakeMarket{},
				StockStorage: &fakeStockStorage{},
				Traffic:      NewStockItemTraffic(),
			},
			args: args{
				msg: domain.Message{
					Transaction: domain.Transaction{
						Slot: domain.Slot{
							Email: "info@example.com",
							StockItem: domain.StockItem{
								MaxPrice: 100,
							},
						},
					},
					Price: 100,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Broker{
				StockStorage:    tt.fields.StockStorage,
				Pricer:          tt.fields.Pricer,
				Market:          tt.fields.Market,
				SMAStack:        tt.fields.SMAStack,
				SettingsStorage: tt.fields.SettingsStorage,
				Stack:           tt.fields.Stack,
				Traffic:         tt.fields.Traffic,
				logger:          tt.fields.logger,
				cron:            tt.fields.cron,
				isShutdown:      tt.fields.isShutdown,
				cronIsRunning:   tt.fields.cronIsRunning,
				isDevMode:       tt.fields.isDevMode,
			}

			done := make(chan struct{})
			var (
				actualSlot []domain.Slot
				actualMsg  []domain.Message
			)
			if len(tt.wantSell) > 0 {
				go func() {
					actualSlot = collectSellMessages(b.Traffic.SellCh)
					done <- struct{}{}
				}()
			}

			if len(tt.wantBuy) > 0 {
				go func() {
					actualMsg = collectMessages(b.Traffic.ConfirmBuyCh)
					done <- struct{}{}
				}()
			}

			if err := b.branchBuyOrSell(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Broker.branchBuyOrSell() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.wantSell) > 0 {
				close(b.Traffic.SellCh)
				<-done
				assert.Equal(t, tt.wantSell, actualSlot)
			}

			if len(tt.wantBuy) > 0 {
				close(b.Traffic.ConfirmBuyCh)
				<-done
				if len(actualMsg) > 0 {
					actualMsg[0].Transaction.Slot.ID = ""
					actualMsg[0].Transaction.BuyAt = time.Time{}
				}
				assert.Equal(t, tt.wantBuy, actualMsg)
			}
		})
	}
}

func collectSellMessages(ch chan domain.Slot) []domain.Slot {
	res := []domain.Slot{}

	for i := range ch {
		res = append(res, i)
	}

	return res
}

func collectMessages(ch chan domain.Message) []domain.Message {
	res := []domain.Message{}

	for i := range ch {
		res = append(res, i)
	}

	return res
}

var errFake = errors.New("fake error")

type fakeSMAStack struct {
	IsGetStackFail bool
	IsTrendUpFail  bool
	TrendUp        bool
	SMAFrame       fakeSMAFrame
}

func (fakeSMAStack) Add(stack string, v float64) bool {
	panic("not implemented")
}

func (this *fakeSMAStack) IsTrendUp(stack string) (bool, error) {
	if this.IsTrendUpFail {
		return false, errFake
	}

	return this.TrendUp, nil
}

func (this *fakeSMAStack) Get(stack string) (domain.SMAFrame, error) {
	if this.IsGetStackFail {
		return nil, errFake
	}

	return &this.SMAFrame, nil
}

func (fakeSMAStack) Reset() {}

func getSMAStack(f fakeSMAStack) domain.SMAStack {
	return &f
}

type fakeSMAFrame struct {
	RangeHigh float64
	RangeLow  float64
	IsEmpty   bool
	LastValue float64
	PrevValue float64
}

func (this *fakeSMAFrame) Distance() float64 { return 0 }
func (this *fakeSMAFrame) IsTrendUp() bool   { return false }
func (this *fakeSMAFrame) Add(v float64)     {}
func (this *fakeSMAFrame) NextCur()          {}
func (this *fakeSMAFrame) CalcAvg()          {}
func (this *fakeSMAFrame) Median() float64   { return 0 }

func (this *fakeSMAFrame) SetRangeHL(h float64, l float64) {}

func (this *fakeSMAFrame) Prev() float64 { return this.PrevValue }
func (this *fakeSMAFrame) Last() float64 { return this.LastValue }

func (this *fakeSMAFrame) RangeHL() (float64, float64) {
	return this.RangeHigh, this.RangeLow
}

func (this *fakeSMAFrame) IsFull() bool {
	return !this.IsEmpty
}

type fakeSettingsStorage struct {
	FakeSettings domain.Settings
	IsFail       bool
}

func (this *fakeSettingsStorage) Settings(context.Context) (domain.Settings, error) {
	if this.IsFail {
		return domain.Settings{}, errFake
	}

	return this.FakeSettings, nil
}

func (fakeSettingsStorage) SaveSettings(context.Context, domain.Settings) error {
	return nil
}

type fakeStack struct {
	GetSlotIsFail bool
	Slots         []domain.Slot
}

func (this *fakeStack) Slot(ctx context.Context, figi string) ([]domain.Slot, error) {
	if this.GetSlotIsFail {
		return nil, errFake
	}

	return this.Slots, nil
}

func (fakeStack) BuyStockItem(_ context.Context, _ domain.Transaction) error {
	return nil
}

func (fakeStack) ConfirmBuyTransaction(_ context.Context, _ domain.Transaction) error {
	return nil
}

func (fakeStack) SellTransaction(_ context.Context, _ domain.Transaction) error {
	return nil
}

type fakeMarket struct{}

func (fakeMarket) ListStockItems(_ context.Context) ([]*domain.StockItem, error) {
	return nil, nil
}

func (fakeMarket) OrderBook(ctx context.Context, i domain.StockItem) (*domain.OrderBook, error) {
	return nil, nil
}

func (fakeMarket) OrderBuy(ctx context.Context, i domain.Transaction) (domain.Transaction, error) {
	return i, nil
}

func (fakeMarket) OrderSell(ctx context.Context, i domain.Transaction) (domain.Transaction, error) {
	return domain.Transaction{}, nil
}

func (fakeMarket) Operations(_ context.Context, _ domain.OperationInput) ([]domain.Transaction, error) {
	return nil, nil
}

type fakeStockStorage struct{}

func (fakeStockStorage) AddStockItemApproved(_ context.Context, _ domain.StockItem) error {
	return nil
}

func (fakeStockStorage) UpdateStockItemApproved(_ context.Context, _ domain.StockItem) error {
	return nil
}

func (fakeStockStorage) RemoveStockItemApproved(_ context.Context, _ domain.StockItem) error {
	return nil
}

func (fakeStockStorage) UpdateActiveStatusStockItemApproved(_ context.Context, _ bool) error {
	return nil
}

func (fakeStockStorage) StockItemApprovedAll(_ context.Context, _ chan domain.Message) {}

func (fakeStockStorage) StockItemApproved(_ context.Context) ([]domain.StockItem, error) {
	return nil, nil
}

func (fakeStockStorage) Slot(_ context.Context, _ string) ([]domain.Slot, error) {
	return nil, nil
}

func (fakeStockStorage) Dealings(_ context.Context) ([]domain.Transaction, error) {
	return nil, nil
}

func (fakeStockStorage) Transaction(_ context.Context, _ string) (domain.Transaction, error) {
	return domain.Transaction{}, nil
}

func (fakeStockStorage) AddMarketPrice(_ context.Context, _ domain.PriceReceiptMessage) error {
	return nil
}

func (fakeStockStorage) Buy(_ context.Context, _ domain.Transaction) error {
	return nil
}

func (fakeStockStorage) ConfirmBuy(_ context.Context, _ domain.Transaction) error {
	return nil
}

func (fakeStockStorage) Sell(_ context.Context, _ domain.Transaction) error {
	return nil
}

func (fakeStockStorage) ConfirmSell(_ context.Context, _ domain.Transaction) error {
	return nil
}

func (fakeStockStorage) PartialSell(_ context.Context, _ domain.Transaction, _ int) error {
	return nil
}

func (fakeStockStorage) PartialConfirmSell(_ context.Context, _ domain.Transaction, _ int) error {
	return nil
}
