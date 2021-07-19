package domain

import "context"

type Stack interface {
	Slot(ctx context.Context, figi string) ([]Slot, error)
	BuyStockItem(context.Context, Transaction) error
	ConfirmBuyTransaction(context.Context, Transaction) error
	SellTransaction(context.Context, Transaction) error
}

type SMAStack interface {
	Add(stack string, v float64) bool
	IsTrendUp(stack string) (bool, error)
	Get(stack string) (SMAFrame, error)
}

type SMAFrame interface {
	Add(v float64)
	NextCur()
	CalcAvg()
	IsTrendUp() bool
	Prev() float64
	IsFull() bool
	Last() float64
	RangeHL() (float64, float64)
	SetRangeHL(h, l float64)
}
