package tinkoff

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

type responseOB struct {
	Payload sdk.RestOrderBook `json:"payload"`
	Status  string            `json:"status"`
}

var errNotNormalTrading = errors.New("trade status isn't normal trading")

func (m *Market) OrderBook(ctx context.Context, i domain.StockItem) (*domain.OrderBook, error) {
	tu, err := extractTokenURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract token from context, %w", err)
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
		return nil, fmt.Errorf("failed to sent request, %w", err)
	}

	if data.Payload.TradeStatus != sdk.NormalTrading {
		return nil, errNotNormalTrading
	}

	bids := make([]domain.PriceQty, len(data.Payload.Bids))
	for i, b := range data.Payload.Bids {
		bids[i] = domain.PriceQty{
			Price: b.Price,
			Qty:   b.Quantity,
		}
	}

	asks := make([]domain.PriceQty, len(data.Payload.Asks))
	for i, b := range data.Payload.Asks {
		asks[i] = domain.PriceQty{
			Price: b.Price,
			Qty:   b.Quantity,
		}
	}

	return &domain.OrderBook{
		StockItem:   i,
		TradeStatus: string(data.Payload.TradeStatus),
		LastPrice:   data.Payload.LastPrice,
		Bids:        bids,
		Asks:        asks,
	}, nil
}
