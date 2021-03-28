package domain

import "context"

type Market interface {
	ListStockItems(context.Context) ([]*StockItem, error)
}
