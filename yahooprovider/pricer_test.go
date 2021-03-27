package yahooprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

func Test_pricer_GetPrice(t *testing.T) {
	type fields struct {
		URL string
	}
	type args struct {
		ctx context.Context
		in  domain.PriceReceiptMessage
	}

	httpwareclient.WithClient(&httpwareclient.HttpClientMock{
		Func: func(req *http.Request) (*http.Response, error) {
			r := response{
				QuoteSummary: quoteSummary{
					Result: []result{
						{
							Price: &price{
								Symbol: "Symbol",
								RegularMarketPrice: priceRaw{
									Raw: 100,
									Fmt: "100",
								},
								MarketState: "REGULAR",
							},
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
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.PriceReceiptMessage
		wantErr bool
	}{
		{
			name: "with regular market",
			args: args{
				ctx: context.Background(),
				in: domain.PriceReceiptMessage{
					Email:     "test@example.com",
					StockItem: domain.StockItem{Ticker: "Symbol"},
				},
			},
			want: domain.PriceReceiptMessage{
				Email:       "test@example.com",
				StockItem:   domain.StockItem{Ticker: "Symbol"},
				Price:       100,
				MarketState: "REGULAR",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pricer{
				URL: tt.fields.URL,
			}
			got, err := p.GetPrice(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("pricer.GetPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pricer.GetPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
