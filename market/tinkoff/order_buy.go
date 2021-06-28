package tinkoff

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

type responseOrderAdd struct {
	Payload sdk.PlacedOrder `json:"payload"`
	Status  string          `json:"status"`
}

type requestOrderAdd struct {
	Lots      int    `json:"lots"`
	Operation string `json:"operation"`
}

func (m *Market) OrderBuy(
	ctx context.Context,
	i domain.Transaction,
) (domain.Transaction, error) {
	tu, err := ExtractTokenURL(ctx)
	if err != nil {
		return domain.Transaction{},
			fmt.Errorf("failed to extract token from context, %w", err)
	}

	if i.FIGI == "" {
		return domain.Transaction{}, fmt.Errorf("FIGI is empty")
	}

	if i.Qty < 1 {
		return domain.Transaction{}, fmt.Errorf("quantity must be more zero")
	}

	dataSend := &requestOrderAdd{
		Lots:      i.Qty,
		Operation: string(sdk.BUY),
	}

	data := &responseOrderAdd{}
	req := &httpwareclient.SendIn{
		Method:   http.MethodPost,
		Headers:  map[string]string{"Authorization": "Bearer " + tu.Token},
		URL:      tu.URL + "/orders/market-order?figi=" + i.FIGI,
		BodySend: dataSend,
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}

	if err := httpwareclient.Send(ctx, req); err != nil {
		return domain.Transaction{},
			fmt.Errorf("failed to sent request, %w", err)
	}

	if data.Status != statusOk {
		return domain.Transaction{}, fmt.Errorf(
			"failed to buy, status is %s, %s",
			data.Status,
			data.Payload.Message,
		)
	}

	result := i
	result.BuyOrderID = data.Payload.ID
	result.Slot.Qty = data.Payload.ExecutedLots

	return result, nil
}
