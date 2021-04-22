package tinkoff

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

func TestMarket_OrderBuy(t *testing.T) {
	type args struct {
		ctx func() context.Context
		i   domain.Slot
	}

	tests := []struct {
		name    string
		args    args
		f       func(req *http.Request) (*http.Response, error)
		want    domain.Slot
		wantErr bool
	}{
		{
			name: "Optimistic",
			args: args{
				ctx: func() context.Context {
					ctx := contexkey.WithAPIURL(context.Background(), "apiurl")
					ctx = contexkey.WithToken(ctx, "token")
					return ctx
				},
				i: domain.Slot{
					Email: "email@example.com",
					StockItem: domain.StockItem{
						FIGI: "figi",
					},
					ID:          "id",
					SlotID:      10,
					StartPrice:  2,
					ChangePrice: 3,
					BuyingPrice: 0,
					Qty:         1,
				},
			},
			f: func(req *http.Request) (*http.Response, error) {
				r := responseOrderAdd{
					Status: statusOk,
					Payload: sdk.PlacedOrder{
						Operation:    sdk.BUY,
						ExecutedLots: 1,
					},
				}
				b, _ := json.Marshal(&r)
				buffer := bytes.NewBuffer(b)
				return &http.Response{
					Body: ioutil.NopCloser(buffer),
				}, nil
			},
			want: domain.Slot{
				Email: "email@example.com",
				StockItem: domain.StockItem{
					FIGI: "figi",
				},
				ID:          "id",
				SlotID:      10,
				StartPrice:  2,
				ChangePrice: 3,
				BuyingPrice: 0,
				Qty:         1,
			},
			wantErr: false,
		},
		{
			name: "without apiurl or token",
			args: args{
				ctx: func() context.Context {
					ctx := contexkey.WithAPIURL(context.Background(), "apiurl")
					ctx = contexkey.WithToken(ctx, "token")
					return ctx
				},
			},
			wantErr: true,
		},
		{
			name: "figi is empty",
			args: args{
				ctx: func() context.Context { return context.Background() },
			},
			wantErr: true,
		},
		{
			name: "quantity is less or equal zero",
			args: args{
				ctx: func() context.Context { return context.Background() },
				i: domain.Slot{
					StockItem: domain.StockItem{
						FIGI: "figi",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "status is not ok",
			args: args{
				ctx: func() context.Context { return context.Background() },
				i: domain.Slot{
					StockItem: domain.StockItem{
						FIGI: "figi",
					},
					Qty: 1,
				},
			},
			f: func(req *http.Request) (*http.Response, error) {
				r := responseOrderAdd{
					Status: statusError,
				}
				b, _ := json.Marshal(&r)
				buffer := bytes.NewBuffer(b)
				return &http.Response{
					Body: ioutil.NopCloser(buffer),
				}, nil
			},
			wantErr: true,
		},
		{
			name: "returns 500",
			args: args{
				ctx: func() context.Context { return context.Background() },
				i: domain.Slot{
					StockItem: domain.StockItem{
						FIGI: "figi",
					},
					Qty: 1,
				},
			},
			f: func(req *http.Request) (*http.Response, error) {
				buffer := bytes.NewBuffer(nil)
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(buffer),
				}, nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Market{}
			httpwareclient.WithClient(&httpwareclient.HttpClientMock{Func: tt.f})
			got, err := m.OrderBuy(tt.args.ctx(), tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.OrderBuy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.OrderBuy() = %v, want %v", got, tt.want)
			}
		})
	}
}
