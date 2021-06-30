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

type responseOrderAdd struct {
	Payload sdk.PlacedOrder `json:"payload"`
	Status  string          `json:"status"`
}

type requestOrderAdd struct {
	Lots      int    `json:"lots"`
	Operation string `json:"operation"`
}

var errBuy = errors.New("failed to buy")

func (m *Market) OrderBuy(
	ctx context.Context,
	i domain.Transaction,
) (domain.Transaction, error) {
	tu, err := extractTokenURL(ctx)
	if err != nil {
		return domain.Transaction{},
			fmt.Errorf("failed to extract token from context, %w", err)
	}

	if i.FIGI == "" {
		return domain.Transaction{}, errFIGIEmpty
	}

	if i.Qty < 1 {
		return domain.Transaction{}, errQuantityZero
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
			"%w, status is %s, %s",
			errBuy, data.Status, data.Payload.Message,
		)
	}

	result := i
	result.BuyOrderID = data.Payload.ID
	result.Slot.Qty = data.Payload.ExecutedLots

	return result, nil
}
