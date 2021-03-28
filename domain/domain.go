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
