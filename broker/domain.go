package broker

type StockItem struct {
	Ticker           string  `json:"ticker"`
	FIGI             string  `json:"figi"`
	AmountLimit      float64 `json:"amount_limit"`
	TransactionLimit int     `json:"transaction_limit"`
}
