package domain

import "context"

type Stack interface {
	Slot(ctx context.Context, figi string) ([]Slot, error)
	BuyStockItem(context.Context, Transaction) error
	ConfirmBuyTransaction(context.Context, Transaction) error
	SellTransaction(context.Context, Transaction) error
}
