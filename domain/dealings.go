package domain

import "time"

type Transaction struct {
	Slot

	SalePrice    float64
	AmountIncome float64

	BuyOrderID  string
	SellOrderID string

	BuyAt    time.Time
	Duration int
	SellAt   time.Time
}
