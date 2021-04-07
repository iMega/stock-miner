package tinkoff

import (
	"context"
	"net/http"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

type responseOB struct {
	Payload sdk.RestOrderBook `json:"payload"`
	Status  string            `json:"status"`
}

func (m *Market) OrderBook(ctx context.Context, i domain.StockItem) (*domain.OrderBook, error) {
	tu, err := ExtractTokenURL(ctx)
	if err != nil {
		return nil, err
	}

	data := &responseOB{}
	req := &httpwareclient.SendIn{
		Method: http.MethodGet,
		Headers: map[string]string{
			"Authorization": "Bearer " + tu.Token,
		},
		URL:      tu.URL + "/market/orderbook?depth=20&figi=" + i.FIGI,
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}

	if err := httpwareclient.Send(ctx, req); err != nil {
		return nil, err
	}

	var bids []domain.PriceQty
	for _, b := range data.Payload.Bids {
		bids = append(bids, domain.PriceQty{
			Price: b.Price,
			Qty:   b.Quantity,
		})
	}

	var asks []domain.PriceQty
	for _, b := range data.Payload.Asks {
		asks = append(asks, domain.PriceQty{
			Price: b.Price,
			Qty:   b.Quantity,
		})
	}

	return &domain.OrderBook{
		StockItem:   i,
		TradeStatus: string(data.Payload.TradeStatus),
		LastPrice:   data.Payload.LastPrice,
		Bids:        bids,
		Asks:        asks,
	}, nil
}
