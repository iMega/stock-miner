package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/imega/stock-miner/domain"
)

func (b *Broker) buy(ctx context.Context, t domain.Transaction) (domain.Transaction, error) {
	tr, err := b.Market.OrderBuy(ctx, t)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to buy item, %s", err)
	}

	if err := b.StockStorage.Buy(ctx, tr); err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to buy item, %#v, %w", tr, err)
	}

	return tr, nil
}

func (b *Broker) confirmBuy(ctx context.Context, t domain.Transaction) error {
	return b.StockStorage.ConfirmBuy(ctx, t)
}

func (b *Broker) confirmBuyJob(tr domain.Transaction) error {
	ctx, err := b.contextWithCreds(context.Background(), tr.Slot.Email)
	if err != nil {
		return fmt.Errorf("failed getting settings, %w", err)
	}

	settings, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting settings, %w", err)
	}

	trs, err := b.Market.Operations(
		ctx,
		domain.OperationInput{
			From:          tr.BuyAt,
			To:            tr.BuyAt.Add(time.Minute),
			OperationType: string(domain.BUY),
			FIGI:          tr.Slot.FIGI,
		},
	)
	if err != nil {
		return fmt.Errorf("failed getting operations, %w", err)
	}

	filteredTR, err := filterOperationByOrderID(trs, tr.BuyOrderID)
	if err != nil {
		return fmt.Errorf("failed to filter operations, %w", err)
	}

	tr.BuyingPrice = filteredTR.BuyingPrice
	tr.Slot.AmountSpent = filteredTR.Slot.AmountSpent
	tr.TargetPrice = calcTargetPrice(
		settings.MarketCommission,
		tr.BuyingPrice,
		settings.GrossMargin,
	)
	tr.Profit = calcSub(tr.TargetPrice, tr.BuyingPrice)

	if err := b.confirmBuy(ctx, tr); err != nil {
		return fmt.Errorf("failed to confirm transaction, %s", err)
	}

	return nil
}
