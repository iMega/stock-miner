package broker

import "context"

type StockStorage interface {
	AddStockItemApproved(context.Context, StockItem) error
	StockItemApproved(context.Context) ([]StockItem, error)
}
