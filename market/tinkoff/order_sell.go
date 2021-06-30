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

var (
	errFIGIEmpty    = errors.New("FIGI is empty")
	errQuantityZero = errors.New("quantity must be more zero")
	errSell         = errors.New("failed to sell")
)

func (m *Market) OrderSell(
	ctx context.Context,
	i domain.Transaction,
) (domain.Transaction, error) {
	tu, err := extractTokenURL(ctx)
	if err != nil {
		return domain.Transaction{},
			fmt.Errorf("failed to extract token from context, %w", err)
	}

	if i.Slot.FIGI == "" {
		return domain.Transaction{}, errFIGIEmpty
	}

	if i.Slot.Qty < 1 {
		return domain.Transaction{}, errQuantityZero
	}

	dataSend := &requestOrderAdd{
		Lots:      i.Qty,
		Operation: string(sdk.SELL),
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
		return domain.Transaction{},
			fmt.Errorf(
				"%w, status: %s, %s",
				errSell, data.Status, data.Payload.Message,
			)
	}

	result := i
	result.SellOrderID = data.Payload.ID
	result.Slot.Qty = data.Payload.ExecutedLots

	return result, nil
}
