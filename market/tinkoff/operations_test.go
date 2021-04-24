package tinkoff

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
	"github.com/stretchr/testify/assert"
)

func Test_maxTradePrices(t *testing.T) {
	type args struct {
		v []sdk.Trade
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "zero trade",
			args: args{
				v: []sdk.Trade{},
			},
			want: 0,
		},
		{
			name: "one trade",
			args: args{
				v: []sdk.Trade{
					{Price: 1},
				},
			},
			want: 1,
		},
		{
			name: "two trade",
			args: args{
				v: []sdk.Trade{
					{Price: 1},
					{Price: 2},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxTradePrices(tt.args.v); got != tt.want {
				t.Errorf("maxTradePrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_Operations(t *testing.T) {
	type args struct {
		ctx func() context.Context
		in  OperationInput
	}
	tests := []struct {
		name    string
		f       func(req *http.Request) (*http.Response, error)
		args    args
		want    []domain.Transaction
		wantErr bool
	}{
		{
			name: "Optimistic",
			f: func(req *http.Request) (*http.Response, error) {
				r := responseOperations{
					Status: statusOk,
					Payload: Operations{
						Operations: []sdk.Operation{
							{
								ID:     "12345678",
								Status: sdk.OperationStatusDone,
								Trades: []sdk.Trade{
									{Price: 1},
									{Price: 1},
									{Price: 2},
								},
								Commission:       sdk.MoneyAmount{Value: 0.5},
								Currency:         sdk.USD,
								Payment:          4,
								Price:            1.75,
								Quantity:         3,
								QuantityExecuted: 3,
								FIGI:             "figi",
								InstrumentType:   sdk.InstrumentTypeStock,
								IsMarginCall:     false,
								DateTime:         time.Now(),
								OperationType:    sdk.BUY,
							},
						},
					},
				}
				b, _ := json.Marshal(&r)
				buffer := bytes.NewBuffer(b)
				return &http.Response{
					Body: ioutil.NopCloser(buffer),
				}, nil
			},
			args: args{
				ctx: func() context.Context {
					ctx := contexkey.WithAPIURL(context.Background(), "apiurl")
					ctx = contexkey.WithToken(ctx, "token")
					return ctx
				},
				in: OperationInput{
					OperationType: string(sdk.BUY),
				},
			},
			want: []domain.Transaction{
				{
					Slot: domain.Slot{
						StockItem: domain.StockItem{
							FIGI: "figi",
						},
						BuyingPrice: 2,
						Qty:         3,
						AmountSpent: 4.5,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Market{}
			httpwareclient.WithClient(&httpwareclient.HttpClientMock{Func: tt.f})
			got, err := m.Operations(tt.args.ctx(), tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Operations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Market.Operations() = %v, want %v", got, tt.want)
			// }
		})
	}
}
