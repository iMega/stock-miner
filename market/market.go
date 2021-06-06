package market

import (
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/market/tinkoff"
)

func New(url, token string) domain.Market {
	return &tinkoff.Market{
		URL:   url,
		Token: token,
	}
}
