package tinkoff

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
	"github.com/imega/stock-miner/money"
)

type Operations struct {
	Operations []sdk.Operation `json:"operations"`
	Message    string          `json:"message,omitempty"`
}

type responseOperations struct {
	Payload Operations `json:"payload"`
	Status  string     `json:"status"`
}

const format = "2006-01-02T15:04:05-07:00"

var errGettingOperations = errors.New("failed getting operations")

func (m *Market) Operations(
	ctx context.Context,
	in domain.OperationInput,
) ([]domain.Transaction, error) {
	tu, err := extractTokenURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from context, %w", err)
	}

	q := url.Values{
		"from": []string{in.From.Format(format)},
		"to":   []string{in.To.Format(format)},
	}
	if in.FIGI != "" {
		q.Add("figi", in.FIGI)
	}

	data := &responseOperations{}
	req := &httpwareclient.SendIn{
		Method:   http.MethodGet,
		Headers:  map[string]string{"Authorization": "Bearer " + tu.Token},
		URL:      tu.URL + "/operations?" + q.Encode(),
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}

	if err := httpwareclient.Send(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to send request, %w", err)
	}

	if data.Status != statusOk {
		return nil,
			fmt.Errorf(
				"%w, status is %s, %s",
				errGettingOperations, data.Status, data.Payload.Message,
			)
	}

	result := make([]domain.Transaction, len(data.Payload.Operations))

	for i, o := range data.Payload.Operations {
		if o.OperationType != sdk.OperationType(in.OperationType) {
			continue
		}

		if o.Status != sdk.OperationStatusDone {
			continue
		}

		if o.Payment == 0 {
			continue
		}

		t := domain.Transaction{
			Slot: domain.Slot{
				StockItem: domain.StockItem{
					FIGI: o.FIGI,
				},
				Qty: o.QuantityExecuted,
			},
		}

		if o.OperationType == sdk.BUY {
			t.BuyOrderID = o.ID
			t.BuyingPrice = maxTradePrices(o.Trades)
			t.Slot.AmountSpent = money.Sum(
				math.Abs(o.Payment),
				math.Abs(o.Commission.Value),
			)
		}

		if o.OperationType == sdk.SELL {
			t.SellOrderID = o.ID
			t.SalePrice = maxTradePrices(o.Trades)
			t.AmountIncome = money.Sub(
				math.Abs(o.Payment),
				math.Abs(o.Commission.Value),
			)
		}

		result[i] = t
	}

	return result, nil
}

func maxTradePrices(v []sdk.Trade) float64 {
	if len(v) == 0 {
		return 0
	}

	tradePrices := make([]float64, len(v))
	for i, trade := range v {
		tradePrices[i] = trade.Price
	}

	sort.Float64s(tradePrices)

	return tradePrices[len(tradePrices)-1]
}
