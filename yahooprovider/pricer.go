package yahooprovider

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/httpwareclient"
)

type Provider struct {
	URL string
}

func New() *Provider {
	return &Provider{
		URL: "https://query1.finance.yahoo.com/v10/finance/quoteSummary/",
	}
}

type response struct {
	QuoteSummary quoteSummary `json:"quoteSummary,omitempty"`
}

type quoteSummary struct {
	Result []result `json:"result,omitempty"`
	Err    *err     `json:"error,omitempty"`
}

type result struct {
	Price *price `json:"price,omitempty"`
}

type price struct {
	Symbol             string   `json:"symbol,omitempty"`
	RegularMarketPrice priceRaw `json:"regularMarketPrice,omitempty"`
	PreMarketPrice     priceRaw `json:"preMarketPrice,omitempty"`
	MarketState        string   `json:"marketState,omitempty"`
}

type priceRaw struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt"`
}

type err struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (p *Provider) Price(
	ctx context.Context,
	i sdk.Instrument,
) (sdk.RestOrderBook, error) {
	data := &response{}

	in := &httpwareclient.SendIn{
		Method:   http.MethodGet,
		URL:      p.URL + i.FIGI + "?modules=price",
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}

	if err := httpwareclient.Send(ctx, in); err != nil {
		return sdk.RestOrderBook{}, err
	}

	if data.QuoteSummary.Err != nil {
		return sdk.RestOrderBook{},
			fmt.Errorf("failed getting price, %s", data.QuoteSummary.Err.Code)
	}

	if len(data.QuoteSummary.Result) == 0 {
		return sdk.RestOrderBook{},
			fmt.Errorf("failed getting price, %s", "empty")
	}

	result := data.QuoteSummary.Result[0]

	price := result.Price.PreMarketPrice.Raw
	if result.Price.MarketState == "REGULAR" {
		price = result.Price.RegularMarketPrice.Raw
	}

	return sdk.RestOrderBook{
		FIGI:      result.Price.Symbol,
		LastPrice: price,
	}, nil
}
