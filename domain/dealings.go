package domain

import "time"

type Transaction struct {
	Slot

	SalePrice    float64 `json:"sale_price"`
	AmountIncome float64 `json:"amount_income"`

	BuyAt    time.Time `json:"buy_at"`
	Duration int       `json:"duration"`
	SellAt   time.Time `json:"sell_at"`
}
