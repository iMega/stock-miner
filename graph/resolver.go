package graph

import "github.com/imega/stock-miner/domain"

type Resolver struct {
	UserStorage      domain.UserStorage
	StockStorage     domain.StockStorage
	MainerController domain.MainerController
}
