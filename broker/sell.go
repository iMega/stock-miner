package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/imega/stock-miner/domain"
)

func (b *Broker) sell(ctx context.Context, t domain.Transaction) (domain.Transaction, error) {
	sellTr, err := b.Market.OrderSell(ctx, t)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to send order sell, %w", err)
	}

	sellTr.SellAt = time.Now()
	if sellTr.Slot.Qty == t.Slot.Qty {
		return sellTr, b.StockStorage.Sell(ctx, sellTr)
	}

	return sellTr, b.StockStorage.PartialSell(ctx, sellTr, t.Slot.Qty)
}

func (b *Broker) confirmSell(ctx context.Context, t domain.Transaction) error {
	in := domain.OperationInput{
		From:          t.SellAt,
		To:            t.SellAt.Add(5 * time.Minute),
		FIGI:          t.Slot.FIGI,
		OperationType: "Sell",
	}
	trs, err := b.Market.Operations(ctx, in)
	if err != nil {
		return fmt.Errorf("failed getting transactions, %w", err)
	}

	filteredTR, err := filterSellOperationByOrderID(trs, t.SellOrderID)
	if err != nil {
		return fmt.Errorf("failed to filter operations, %w", err)
	}

	t.SalePrice = filteredTR.SalePrice
	t.AmountIncome = filteredTR.AmountIncome
	t.Duration = int(t.SellAt.Unix() - t.BuyAt.Unix())
	t.TotalProfit = calcSub(t.AmountIncome, t.AmountSpent)

	if t.Slot.Qty == filteredTR.Qty {
		return b.StockStorage.ConfirmSell(ctx, t)
	}

	return b.StockStorage.PartialConfirmSell(ctx, t, filteredTR.Qty)
}

func filterSellOperationByOrderID(trs []domain.Transaction, orderID string) (domain.Transaction, error) {
	for _, t := range trs {
		if t.SellOrderID == orderID {
			return t, nil
		}
	}

	return domain.Transaction{}, fmt.Errorf("operation does not exist")
}

func (b *Broker) confirmSellJob(msg domain.Message) error {
	ctx, err := b.contextWithCreds(
		context.Background(),
		msg.Transaction.Slot.Email,
	)
	if err != nil {
		return fmt.Errorf("failed getting creds, %w", err)
	}

	if err := b.confirmSell(ctx, msg.Transaction); err != nil {
		return fmt.Errorf("failed to confirm sell, %w", err)
	}

	return nil
}
