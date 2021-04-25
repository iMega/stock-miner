package domain

import (
	"context"
	"time"
)

type Market interface {
	ListStockItems(context.Context) ([]*StockItem, error)
	OrderBook(ctx context.Context, i StockItem) (*OrderBook, error)
	OrderBuy(ctx context.Context, i Transaction) (Transaction, error)
	Operations(context.Context, OperationInput) ([]Transaction, error)
}

type OperationInput struct {
	From          time.Time `json:"from"`
	To            time.Time `json:"to"`
	FIGI          string    `json:"figi"`
	OperationType string    `json:"operation_type"`
}
