// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Deal struct {
	ID           string   `json:"id"`
	Ticker       string   `json:"ticker"`
	Figi         string   `json:"figi"`
	StartPrice   float64  `json:"startPrice"`
	ChangePrice  float64  `json:"changePrice"`
	BuyingPrice  *float64 `json:"buyingPrice"`
	TargetPrice  *float64 `json:"targetPrice"`
	Profit       *float64 `json:"profit"`
	SalePrice    *float64 `json:"salePrice"`
	Qty          *int     `json:"qty"`
	AmountSpent  *float64 `json:"amountSpent"`
	AmountIncome *float64 `json:"amountIncome"`
	TotalProfit  *float64 `json:"totalProfit"`
	BuyAt        *string  `json:"buyAt"`
	Duration     *int     `json:"duration"`
	SellAt       *string  `json:"sellAt"`
	Currency     string   `json:"currency"`
}

type MarketCredentials struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	APIURL string `json:"apiUrl"`
}

type MarketCredentialsInput struct {
	Name   string `json:"name"`
	APIURL string `json:"apiUrl"`
	Token  string `json:"token"`
}

type MemStats struct {
	Alloc      string `json:"alloc"`
	TotalAlloc string `json:"totalAlloc"`
	Sys        string `json:"sys"`
}

type RulePriceInput struct {
	MarketCommission *float64 `json:"marketCommission"`
	GrossMargin      *float64 `json:"grossMargin"`
}

type Settings struct {
	Slot              *SlotSettings        `json:"slot"`
	MarketCredentials []*MarketCredentials `json:"marketCredentials"`
	MarketProvider    string               `json:"marketProvider"`
	MarketCommission  *float64             `json:"marketCommission"`
	GrossMargin       *float64             `json:"grossMargin"`
	MiningStatus      bool                 `json:"miningStatus"`
}

type Slot struct {
	ID           string   `json:"id"`
	Ticker       string   `json:"ticker"`
	Figi         string   `json:"figi"`
	StartPrice   float64  `json:"startPrice"`
	ChangePrice  float64  `json:"changePrice"`
	BuyingPrice  *float64 `json:"buyingPrice"`
	TargetPrice  *float64 `json:"targetPrice"`
	Profit       *float64 `json:"profit"`
	Qty          *int     `json:"qty"`
	AmountSpent  *float64 `json:"amountSpent"`
	TargetAmount *float64 `json:"targetAmount"`
	TotalProfit  *float64 `json:"totalProfit"`
	Currency     string   `json:"currency"`
	CurrentPrice float64  `json:"currentPrice"`
}

type SlotSettings struct {
	Volume              int      `json:"volume"`
	ModificatorMinPrice *float64 `json:"modificatorMinPrice"`
}

type SlotSettingsInput struct {
	Volume              int      `json:"volume"`
	ModificatorMinPrice *float64 `json:"modificatorMinPrice"`
}

type StockItem struct {
	Ticker            string   `json:"ticker"`
	Figi              string   `json:"figi"`
	Isin              *string  `json:"isin"`
	MinPriceIncrement *float64 `json:"minPriceIncrement"`
	Lot               *int     `json:"lot"`
	Currency          *string  `json:"currency"`
	Name              *string  `json:"name"`
	AmountLimit       float64  `json:"amountLimit"`
	TransactionLimit  int      `json:"transactionLimit"`
	StartTime         int      `json:"startTime"`
	EndTime           int      `json:"endTime"`
	Active            bool     `json:"active"`
	MaxPrice          float64  `json:"maxPrice"`
}

type StockItemInput struct {
	Ticker           string  `json:"ticker"`
	Figi             string  `json:"figi"`
	AmountLimit      float64 `json:"amountLimit"`
	TransactionLimit int     `json:"transactionLimit"`
	Currency         string  `json:"currency"`
	StartTime        int     `json:"startTime"`
	EndTime          int     `json:"endTime"`
	Active           bool    `json:"active"`
	MaxPrice         float64 `json:"maxPrice"`
}

type User struct {
	Email  string  `json:"email"`
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
	Role   *string `json:"role"`
}

type UserInput struct {
	Email  string  `json:"email"`
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
}
