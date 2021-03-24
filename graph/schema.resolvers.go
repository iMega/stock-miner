package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/imega/stock-miner/broker"
	"github.com/imega/stock-miner/graph/generated"
	"github.com/imega/stock-miner/graph/model"
)

func (r *mutationResolver) AddStockItemApproved(ctx context.Context, item model.StockItemInput) (*model.StockItem, error) {
	in := broker.StockItem{
		Ticker:           item.Ticker,
		FIGI:             item.Figi,
		AmountLimit:      item.AmountLimit,
		TransactionLimit: item.TransactionLimit,
	}
	if err := r.StockStorage.AddStockItemApproved(ctx, in); err != nil {
		return nil, err
	}

	return &model.StockItem{
		Ticker:           item.Ticker,
		Figi:             item.Figi,
		AmountLimit:      item.AmountLimit,
		TransactionLimit: item.TransactionLimit,
	}, nil
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
