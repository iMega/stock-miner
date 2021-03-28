package market

import (
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/market/tinkoff"
)

func New(URL, token string) domain.Market {
	return &tinkoff.Market{
		URL:   URL,
		Token: token,
	}
}
