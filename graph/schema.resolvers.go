package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/graph/generated"
	"github.com/imega/stock-miner/graph/model"
	"github.com/imega/stock-miner/stats"
)

func (r *mutationResolver) AddStockItemApproved(ctx context.Context, items []*model.StockItemInput) (bool, error) {
	for _, item := range items {
		in := domain.StockItem{
			Ticker:           item.Ticker,
			FIGI:             item.Figi,
			AmountLimit:      item.AmountLimit,
			TransactionLimit: item.TransactionLimit,
		}
		if err := r.StockStorage.AddStockItemApproved(ctx, in); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *mutationResolver) GlobalMiningStop(ctx context.Context) (bool, error) {
	return r.MainerController.Stop(), nil
}

func (r *mutationResolver) GlobalMiningStart(ctx context.Context) (bool, error) {
	return r.MainerController.Start(), nil
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	user, err := r.UserStorage.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Email:  user.Email,
		Name:   &user.Name,
		Avatar: &user.Avatar,
	}, nil
}

func (r *queryResolver) StockItemApproved(ctx context.Context) ([]*model.StockItem, error) {
	var result []*model.StockItem
	items, err := r.StockStorage.StockItemApproved(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		result = append(result, &model.StockItem{
			Ticker:           item.Ticker,
			Figi:             item.FIGI,
			AmountLimit:      item.AmountLimit,
			TransactionLimit: item.TransactionLimit,
		})
	}

	return result, nil
}

func (r *queryResolver) MemStats(ctx context.Context) (*model.MemStats, error) {
	m := stats.GetMemStats()

	return &model.MemStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
	}, nil
}

func (r *queryResolver) GlobalMiningStatus(ctx context.Context) (bool, error) {
	return r.MainerController.Status(), nil
}

func (r *queryResolver) MarketStockItems(ctx context.Context) ([]*model.StockItem, error) {
	items, err := r.Market.ListStockItems(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.StockItem
	for _, item := range items {
		result = append(result, &model.StockItem{
			Ticker:            item.Ticker,
			Figi:              item.FIGI,
			Isin:              &item.ISIN,
			MinPriceIncrement: &item.MinPriceIncrement,
			Lot:               &item.Lot,
			Currency:          &item.Currency,
			Name:              &item.Name,
		})
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
