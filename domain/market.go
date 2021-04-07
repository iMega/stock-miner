package domain

import "context"

type Market interface {
	ListStockItems(context.Context) ([]*StockItem, error)
	OrderBook(ctx context.Context, i StockItem) (*OrderBook, error)
}
