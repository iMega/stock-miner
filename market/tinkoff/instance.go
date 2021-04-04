package tinkoff

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

type Market struct {
	URL   string
	Token string
}

type response struct {
	Payload payload `json:"payload"`
	Status  string  `json:"status"`
}

type payload struct {
	Total       int64            `json:"total"`
	Instruments []sdk.Instrument `json:"instruments"`
}

func (m *Market) ListStockItems(ctx context.Context) ([]*domain.StockItem, error) {
	token, ok := contexkey.TokenFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract token from context")
	}

	apiurl, ok := contexkey.APIURLFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract apiurl from context")
	}

	data := &response{}
	req := &httpwareclient.SendIn{
		Method: http.MethodGet,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		URL:      apiurl + "/market/stocks",
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}

	if err := httpwareclient.Send(ctx, req); err != nil {
		return nil, err
	}

	var result []*domain.StockItem
	for _, i := range data.Payload.Instruments {
		result = append(result, &domain.StockItem{
			Ticker:            i.Ticker,
			FIGI:              i.FIGI,
			ISIN:              i.ISIN,
			Name:              i.Name,
			MinPriceIncrement: i.MinPriceIncrement,
			Lot:               i.Lot,
			Currency:          string(i.Currency),
		})
	}

	return result, nil
}
