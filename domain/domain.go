package domain

type StockItem struct {
	Ticker            string  `json:"ticker"`
	FIGI              string  `json:"figi"`
	AmountLimit       float64 `json:"amount_limit"`
	TransactionLimit  int     `json:"transaction_limit"`
	ISIN              string  `json:"isin"`
	Name              string  `json:"name"`
	MinPriceIncrement float64 `json:"minPriceIncrement"`
	Lot               uint8   `json:"lot"`
	Currency          string  `json:"currency"`
	StartTime         uint8
	EndTime           uint8
	IsActive          bool
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

	Currency string `json:"currency"`
}

type OrderBook struct {
	StockItem
	Bids        []PriceQty `json:"bids"`
	Asks        []PriceQty `json:"asks"`
	TradeStatus string     `json:"tradeStatus"`
	LastPrice   float64    `json:"lastPrice"`
}

type PriceQty struct {
	Price float64
	Qty   float64
}

type OperationType string

const (
	// BUY operation type.
	BUY OperationType = "Buy"

	// SELL operation type.
	SELL OperationType = "Sell"
)

type TaskOperation struct {
	Attempt     int
	Transaction Transaction
	Operation   OperationType
}

type Message struct {
	RetryCount  int
	Error       error
	Price       float64
	Transaction Transaction
}
