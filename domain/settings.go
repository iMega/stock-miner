package domain

import "context"

type SettingsStorage interface {
	Settings(context.Context) (Settings, error)
	SaveSettings(context.Context, Settings) error
}

type Settings struct {
	Slot SlotSettings `json:"slot,omitempty"`

	MarketCredentials map[string]MarketCredentials `json:"market_credentials,omitempty"`
}

type SlotSettings struct {
	Volume int `json:"volume,omitempty"`
}

type MarketCredentials struct {
	Token  string `json:"token,omitempty"`
	APIURL string `json:"api_url,omitempty"`
}
