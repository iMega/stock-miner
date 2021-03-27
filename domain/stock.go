package domain

import "context"

type StockStorage interface {
	AddStockItemApproved(context.Context, StockItem) error
	StockItemApprovedAll(context.Context, chan PriceReceiptMessage)
	StockItemApproved(context.Context) ([]StockItem, error)
}
