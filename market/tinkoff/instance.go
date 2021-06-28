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

const (
	statusOk = "Ok"
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

type tokenURL struct {
	Token string
	URL   string
}

func ExtractTokenURL(ctx context.Context) (*tokenURL, error) {
	token, ok := contexkey.TokenFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract token from context")
	}

	apiurl, ok := contexkey.APIURLFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract apiurl from context")
	}

	return &tokenURL{
		Token: token,
		URL:   apiurl,
	}, nil
}

func (m *Market) ListStockItems(ctx context.Context) ([]*domain.StockItem, error) {
	tu, err := ExtractTokenURL(ctx)
	if err != nil {
		return nil, err
	}

	data := &response{}
	req := &httpwareclient.SendIn{
		Method: http.MethodGet,
		Headers: map[string]string{
			"Authorization": "Bearer " + tu.Token,
		},
		URL:      tu.URL + "/market/stocks",
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}

	if err := httpwareclient.Send(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to sent request, %w", err)
	}

	result := make([]*domain.StockItem, len(data.Payload.Instruments))
	for idx, i := range data.Payload.Instruments {
		result[idx] = &domain.StockItem{
			Ticker:            i.Ticker,
			FIGI:              i.FIGI,
			ISIN:              i.ISIN,
			Name:              i.Name,
			MinPriceIncrement: i.MinPriceIncrement,
			Lot:               i.Lot,
			Currency:          string(i.Currency),
		}
	}

	return result, nil
}
