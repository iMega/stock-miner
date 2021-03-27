package domain

type StockItem struct {
	Ticker           string  `json:"ticker"`
	FIGI             string  `json:"figi"`
	AmountLimit      float64 `json:"amount_limit"`
	TransactionLimit int     `json:"transaction_limit"`
}

type PriceReceiptMessage struct {
	Email       string
	Price       float64
	MarketState string
	StockItem
	Error error
}
