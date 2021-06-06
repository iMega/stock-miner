package yahooprovider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

// URL: https://query1.finance.yahoo.com/v10/finance/quoteSummary/

type pricer struct {
	URL string
}

func New(url string) domain.Pricer {
	p := &pricer{URL: url}

	return p
}

func (p *pricer) GetPrice(
	ctx context.Context,
	in domain.PriceReceiptMessage,
) (domain.PriceReceiptMessage, error) {
	data := &response{}
	result := domain.PriceReceiptMessage{
		StockItem: in.StockItem,
		Email:     in.Email,
	}

	req := &httpwareclient.SendIn{
		Method:   http.MethodGet,
		URL:      p.URL + in.Ticker + "?modules=price",
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}

	if err := httpwareclient.Send(ctx, req); err != nil {
		return result, err
	}

	if data.QuoteSummary.Err != nil {
		return result, fmt.Errorf(
			"failed getting price, %s",
			data.QuoteSummary.Err.Code,
		)
	}

	if len(data.QuoteSummary.Result) == 0 {
		return result, fmt.Errorf("failed getting price, %s", "empty")
	}

	response := data.QuoteSummary.Result[0]

	price := response.Price.PreMarketPrice.Raw
	if response.Price.MarketState == "REGULAR" {
		price = response.Price.RegularMarketPrice.Raw
	}

	result.Price = price
	result.MarketState = response.Price.MarketState

	return result, nil
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
