package domain

type StockItem struct {
	Ticker            string  `json:"ticker"`
	FIGI              string  `json:"figi"`
	AmountLimit       float64 `json:"amount_limit"`
	TransactionLimit  int     `json:"transaction_limit"`
	ISIN              string  `json:"isin"`
	Name              string  `json:"name"`
	MinPriceIncrement float64 `json:"minPriceIncrement"`
	Lot               int     `json:"lot"`
	Currency          string  `json:"currency"`
}

type PriceReceiptMessage struct {
	Email       string
	Price       float64
	MarketState string
	StockItem
	Error error
}

type Slot struct {
	Email string
	StockItem

	ID     string `json:"id"`
	SlotID int    `json:"slot_id"`

	StartPrice  float64 `json:"start_price"`
	ChangePrice float64 `json:"change_price"`
	BuyingPrice float64 `json:"buying_price"`
	TargetPrice float64 `json:"target_price"`
	Profit      float64 `json:"profit"`

	Qty int `json:"qty"`

	TargetAmount float64 `json:"target_amount"`
	AmountSpent  float64 `json:"amount_spent"`
	TotalProfit  float64 `json:"total_profit"`
}
