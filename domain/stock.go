package domain

import "context"

type StockStorage interface {
	AddStockItemApproved(context.Context, StockItem) error
	UpdateStockItemApproved(context.Context, StockItem) error
	RemoveStockItemApproved(context.Context, StockItem) error

	StockItemApprovedAll(context.Context, chan PriceReceiptMessage)
	StockItemApproved(context.Context) ([]StockItem, error)
	Slot(context.Context, string) ([]Slot, error)
	Dealings(context.Context) ([]Transaction, error)

	AddMarketPrice(context.Context, PriceReceiptMessage) error

	Buy(context.Context, Transaction) error
	ConfirmBuy(context.Context, Transaction) error
}
