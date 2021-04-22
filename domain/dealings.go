package domain

import "time"

type Transaction struct {
	Slot

	SalePrice    float64 `json:"sale_price"`
	AmountIncome float64 `json:"amount_income"`

	BuyOrderID  string `json:"buy_order_id"`
	SellOrderID string `json:"sell_order_id"`

	BuyAt    time.Time `json:"buy_at"`
	Duration int       `json:"duration"`
	SellAt   time.Time `json:"sell_at"`
}
