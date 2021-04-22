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

func (m *Market) OrderBuy(ctx context.Context, i domain.Slot) (domain.Slot, error) {
	tu, err := ExtractTokenURL(ctx)
	if err != nil {
		return domain.Slot{}, err
	}

	if i.FIGI == "" {
		return domain.Slot{}, fmt.Errorf("FIGI is empty")
	}

	if i.Qty < 1 {
		return domain.Slot{}, fmt.Errorf("quantity must be more zero")
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
		return domain.Slot{}, err
	}

	if data.Status != statusOk {
		return domain.Slot{}, fmt.Errorf(
			"failed to buy, status is %s, %s",
			data.Status,
			data.Payload.Message,
		)
	}

	return domain.Slot{
		Email:       i.Email,
		StockItem:   i.StockItem,
		ID:          i.ID,
		SlotID:      i.SlotID,
		StartPrice:  i.StartPrice,
		ChangePrice: i.ChangePrice,
		BuyingPrice: 0,
		Qty:         data.Payload.ExecutedLots,
	}, nil
}
