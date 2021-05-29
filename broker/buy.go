package broker

import (
	"context"
	"fmt"

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
