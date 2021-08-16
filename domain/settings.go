package domain

import "context"

type SettingsStorage interface {
	Settings(context.Context) (Settings, error)
	SaveSettings(context.Context, Settings) error
}

type Settings struct {
	Slot SlotSettings `json:"slot,omitempty"`

	MarketCredentials map[string]MarketCredentials `json:"market_credentials,omitempty"`
	MarketProvider    string                       `json:"market_provider,omitempty"`
	MarketCommission  float64                      `json:"market_commission,omitempty"`

	GrossMargin float64 `json:"gross_margin,omitempty"`

	MiningStatus bool `json:"miningStatus"`

	MainSettings MainSettings
}

type SlotSettings struct {
	Volume              int     `json:"volume,omitempty"`
	ModificatorMinPrice float64 `json:"modificator_min_price,omitempty"`
}

type MarketCredentials struct {
	Token  string `json:"token,omitempty"`
	APIURL string `json:"api_url,omitempty"`
}

type MainSettings struct {
	MiningStatus bool
}
