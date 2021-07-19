package domain

type StockItem struct {
	Ticker            string
	FIGI              string
	ISIN              string
	Name              string
	Currency          string
	AmountLimit       float64
	MinPriceIncrement float64
	TransactionLimit  int
	Lot               uint8
	StartTime         uint8
	EndTime           uint8
	IsActive          bool
}

type StockItemRange struct {
	High float64
	Low  float64
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

	ID     string
	SlotID int

	StartPrice  float64
	ChangePrice float64
	BuyingPrice float64
	TargetPrice float64
	Profit      float64

	Qty int

	TargetAmount float64
	AmountSpent  float64
	TotalProfit  float64

	Currency string
}

type OrderBook struct {
	StockItem
	Bids        []PriceQty
	Asks        []PriceQty
	TradeStatus string
	LastPrice   float64
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
