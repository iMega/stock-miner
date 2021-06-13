package tinkoff

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
	"github.com/stretchr/testify/assert"
)

func TestMarket_OrderBook(t *testing.T) {
	type args struct {
		ctx func() context.Context
		i   domain.StockItem
	}

	tests := []struct {
		name    string
		args    args
		f       func(req *http.Request) (*http.Response, error)
		want    *domain.OrderBook
		wantErr bool
	}{
		{
			name: "trading status is OpeningPeriod",
			f: func(req *http.Request) (*http.Response, error) {
				r := responseOB{
					Status: statusOk,
					Payload: sdk.RestOrderBook{
						TradeStatus: sdk.OpeningPeriod,
					},
				}

				b, err := json.Marshal(&r)
				if err != nil {
					return nil, err
				}

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
				i: domain.StockItem{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Market{}
			httpwareclient.WithClient(&httpwareclient.HTTPClientMock{Func: tt.f})
			got, err := m.OrderBook(tt.args.ctx(), tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.OrderBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
