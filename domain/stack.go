package domain

import "context"

type Stack interface {
	Slot(ctx context.Context, figi string) ([]Slot, error)
	BuyStockItem(context.Context, Slot) error
}
