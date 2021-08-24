package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/graph/generated"
	"github.com/imega/stock-miner/graph/model"
	"github.com/imega/stock-miner/money"
	"github.com/imega/stock-miner/stats"
)

func (r *mutationResolver) AddStockItemApproved(ctx context.Context, items []*model.StockItemInput) (bool, error) {
	settings, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting settings, %w", err)
	}

	for _, item := range items {
		in := domain.StockItem{
			Ticker:           item.Ticker,
			FIGI:             item.Figi,
			AmountLimit:      item.AmountLimit,
			TransactionLimit: item.TransactionLimit,
			Currency:         item.Currency,
			StartTime:        uint8(item.StartTime),
			EndTime:          uint8(item.EndTime),
			IsActive:         settings.MiningStatus,
		}
		if err := r.StockStorage.AddStockItemApproved(ctx, in); err != nil {
			return false,
				fmt.Errorf("failed to add approved stock item, %w", err)
		}
	}

	return true, nil
}

func (r *mutationResolver) RemoveStockItemApproved(ctx context.Context, items []*model.StockItemInput) (bool, error) {
	for _, item := range items {
		in := domain.StockItem{
			Ticker:           item.Ticker,
			FIGI:             item.Figi,
			AmountLimit:      item.AmountLimit,
			TransactionLimit: item.TransactionLimit,
			Currency:         item.Currency,
		}
		if err := r.StockStorage.RemoveStockItemApproved(ctx, in); err != nil {
			return false,
				fmt.Errorf("failed to remove approved stock item, %w", err)
		}
	}

	return true, nil
}

func (r *mutationResolver) UpdateStockItemApproved(ctx context.Context, items []*model.StockItemInput) (bool, error) {
	for _, item := range items {
		in := domain.StockItem{
			Ticker:           item.Ticker,
			FIGI:             item.Figi,
			AmountLimit:      item.AmountLimit,
			TransactionLimit: item.TransactionLimit,
			Currency:         item.Currency,
			StartTime:        uint8(item.StartTime),
			EndTime:          uint8(item.EndTime),
		}
		if err := r.StockStorage.UpdateStockItemApproved(ctx, in); err != nil {
			return false,
				fmt.Errorf("failed to update approved stock item, %w", err)
		}
	}

	return true, nil
}

func (r *mutationResolver) EnableStockItemsApproved(ctx context.Context) (bool, error) {
	err := r.StockStorage.UpdateActiveStatusStockItemApproved(ctx, true)
	if err != nil {
		return false,
			fmt.Errorf("failed to enable approved stock items, %w", err)
	}

	settings, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting settings, %w", err)
	}

	settings.MiningStatus = true

	if err := r.SettingsStorage.SaveSettings(ctx, settings); err != nil {
		return false, fmt.Errorf("failed to save settings, %w", err)
	}

	return true, nil
}

func (r *mutationResolver) DisableStockItemsApproved(ctx context.Context) (bool, error) {
	err := r.StockStorage.UpdateActiveStatusStockItemApproved(ctx, false)
	if err != nil {
		return false,
			fmt.Errorf("failed to disable approved stock items, %w", err)
	}

	settings, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting settings, %w", err)
	}

	settings.MiningStatus = false

	if err := r.SettingsStorage.SaveSettings(ctx, settings); err != nil {
		return false, fmt.Errorf("failed to save settings, %w", err)
	}

	return true, nil
}

func (r *mutationResolver) MarketCredentials(ctx context.Context, creds model.MarketCredentialsInput) (bool, error) {
	s, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to save creds, %w", err)
	}

	if _, ok := s.MarketCredentials[creds.Name]; !ok {
		if s.MarketCredentials == nil {
			s.MarketCredentials = make(map[string]domain.MarketCredentials)
		}

		s.MarketCredentials[creds.Name] = domain.MarketCredentials{}
	}

	s.MarketCredentials[creds.Name] = domain.MarketCredentials{
		Token:  creds.Token,
		APIURL: creds.APIURL,
	}
	s.MarketProvider = creds.Name

	if err := r.SettingsStorage.SaveSettings(ctx, s); err != nil {
		return false, fmt.Errorf("failed to save creds, %w", err)
	}

	return true, nil
}

func (r *mutationResolver) Slot(ctx context.Context, global model.SlotSettingsInput) (bool, error) {
	s, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting settings, %w", err)
	}

	s.Slot.Volume = global.Volume
	if global.ModificatorMinPrice != nil {
		s.Slot.ModificatorMinPrice = *global.ModificatorMinPrice
	}

	if err := r.SettingsStorage.SaveSettings(ctx, s); err != nil {
		return false, fmt.Errorf("failed to save slot, %w", err)
	}

	return true, nil
}

func (r *mutationResolver) RulePrice(ctx context.Context, global model.RulePriceInput) (bool, error) {
	s, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting settings, %w", err)
	}

	s.MarketCommission = *global.MarketCommission
	s.GrossMargin = *global.GrossMargin

	if err := r.SettingsStorage.SaveSettings(ctx, s); err != nil {
		return false, fmt.Errorf("failed to save rule price, %w", err)
	}

	return true, nil
}

func (r *mutationResolver) GlobalMiningStop(ctx context.Context) (bool, error) {
	u, err := r.UserStorage.GetUser(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting user, %w", err)
	}

	if u.Role != "root" {
		return false, nil
	}

	return r.MainerController.Stop(), nil
}

func (r *mutationResolver) GlobalMiningStart(ctx context.Context) (bool, error) {
	u, err := r.UserStorage.GetUser(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting user, %w", err)
	}

	if u.Role != "root" {
		return false, nil
	}

	return r.MainerController.Start(), nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, user model.UserInput) (bool, error) {
	u, err := r.UserStorage.GetUser(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting user, %w", err)
	}

	if u.Role != "root" {
		return false, nil
	}

	err = r.UserStorage.CreateUser(
		ctx,
		domain.User{
			Email: user.Email,
			Role:  "user",
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to create user, %w", err)
	}

	return true, nil
}

func (r *mutationResolver) RemoveUser(ctx context.Context, user model.UserInput) (bool, error) {
	u, err := r.UserStorage.GetUser(ctx)
	if err != nil {
		return false, fmt.Errorf("failed getting user, %w", err)
	}

	if u.Role != "root" {
		return false, nil
	}

	err = r.UserStorage.RemoveUser(ctx, domain.User{Email: user.Email})
	if err != nil {
		return false, fmt.Errorf("failed to remove user, %w", err)
	}

	return true, nil
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	user, err := r.UserStorage.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting user, %w", err)
	}

	return &model.User{
		Email:  user.Email,
		Name:   &user.Name,
		Avatar: &user.Avatar,
		Role:   &user.Role,
	}, nil
}

func (r *queryResolver) StockItemApproved(ctx context.Context) ([]*model.StockItem, error) {
	items, err := r.StockStorage.StockItemApproved(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting approved stock item, %w", err)
	}

	result := make([]*model.StockItem, len(items))
	for i, item := range items {
		result[i] = &model.StockItem{
			Ticker:           item.Ticker,
			Figi:             item.FIGI,
			AmountLimit:      item.AmountLimit,
			TransactionLimit: item.TransactionLimit,
			Currency:         &item.Currency,
			StartTime:        int(item.StartTime),
			EndTime:          int(item.EndTime),
		}
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
	s, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting settings, %w", err)
	}

	ctxNew := contexkey.WithToken(
		ctx,
		s.MarketCredentials[s.MarketProvider].Token,
	)
	ctxNew = contexkey.WithAPIURL(
		ctxNew,
		s.MarketCredentials[s.MarketProvider].APIURL,
	)

	items, err := r.Market.ListStockItems(ctxNew)
	if err != nil {
		return nil, fmt.Errorf("failed getting stock items, %w", err)
	}

	result := make([]*model.StockItem, len(items))

	for i, item := range items {
		lot := int(item.Lot)
		result[i] = &model.StockItem{
			Ticker:            item.Ticker,
			Figi:              item.FIGI,
			Isin:              &item.ISIN,
			MinPriceIncrement: &item.MinPriceIncrement,
			Lot:               &lot,
			Currency:          &item.Currency,
			Name:              &item.Name,
		}
	}

	return result, nil
}

func (r *queryResolver) Settings(ctx context.Context) (*model.Settings, error) {
	s, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting settings, %w", err)
	}

	cred := []*model.MarketCredentials{}
	for k, c := range s.MarketCredentials {
		cred = append(cred, &model.MarketCredentials{
			Name:   k,
			Token:  c.Token,
			APIURL: c.APIURL,
		})
	}

	return &model.Settings{
		Slot: &model.SlotSettings{
			Volume:              s.Slot.Volume,
			ModificatorMinPrice: &s.Slot.ModificatorMinPrice,
		},
		MarketCredentials: cred,
		MarketCommission:  &s.MarketCommission,
		GrossMargin:       &s.GrossMargin,
		MiningStatus:      s.MiningStatus,
	}, nil
}

func (r *queryResolver) Slots(ctx context.Context) ([]*model.Slot, error) {
	s, err := r.SettingsStorage.Settings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting settings, %w", err)
	}

	slots, err := r.StockStorage.Slot(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed getting slot, %w", err)
	}

	result := make([]*model.Slot, len(slots))

	for i, v := range slots {
		slot := v

		var (
			currentPrice float64
			currency     = "USD"
		)

		if frame, err := r.SMAStack.Get(slot.Ticker); err == nil {
			currentPrice = frame.Last()
		}

		if len(slot.StockItem.Currency) > 0 {
			currency = slot.StockItem.Currency
		}

		p := money.Procent(slot.TargetAmount, s.MarketCommission)
		profit := money.Sub(slot.TargetAmount, money.Sum(slot.AmountSpent, p))

		result[i] = &model.Slot{
			ID:           slot.ID,
			Ticker:       slot.Ticker,
			Figi:         slot.FIGI,
			StartPrice:   slot.StartPrice,
			ChangePrice:  slot.ChangePrice,
			BuyingPrice:  &slot.BuyingPrice,
			TargetPrice:  &slot.TargetPrice,
			Profit:       &slot.Profit,
			Qty:          &slot.Qty,
			AmountSpent:  &slot.AmountSpent,
			TargetAmount: &slot.TargetAmount,
			TotalProfit:  &profit,
			Currency:     currency,
			CurrentPrice: currentPrice,
		}
	}

	return result, nil
}

func (r *queryResolver) Dealings(ctx context.Context) ([]*model.Deal, error) {
	dealings, err := r.StockStorage.Dealings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting dealings, %w", err)
	}

	result := make([]*model.Deal, len(dealings))

	for i, v := range dealings {
		deal := v
		buyAt := deal.BuyAt.Format(time.RFC3339)
		sellAt := deal.SellAt.Format(time.RFC3339)

		result[i] = &model.Deal{
			ID:           deal.ID,
			Ticker:       deal.Ticker,
			Figi:         deal.FIGI,
			StartPrice:   deal.Slot.StartPrice,
			ChangePrice:  deal.Slot.ChangePrice,
			BuyingPrice:  &deal.Slot.BuyingPrice,
			TargetPrice:  &deal.Slot.TargetPrice,
			Profit:       &deal.Slot.Profit,
			SalePrice:    &deal.SalePrice,
			Qty:          &deal.Slot.Qty,
			AmountSpent:  &deal.Slot.AmountSpent,
			AmountIncome: &deal.AmountIncome,
			TotalProfit:  &deal.TotalProfit,
			BuyAt:        &buyAt,
			Duration:     &deal.Duration,
			SellAt:       &sellAt,
			Currency:     "USD",
		}
	}

	return result, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	user, err := r.UserStorage.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting user, %w", err)
	}

	if user.Role != "root" {
		return nil, nil
	}

	users, err := r.UserStorage.Users(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting users, %w", err)
	}

	result := make([]*model.User, len(users))
	for i, v := range users {
		user := v
		result[i] = &model.User{
			Email:  user.Email,
			Name:   &user.Name,
			Avatar: &user.Avatar,
		}
	}

	return result, nil
}

func (r *subscriptionResolver) MemStats(ctx context.Context) (<-chan *model.MemStats, error) {
	ch := make(chan *model.MemStats)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)

				return

			case <-time.After(1 * time.Second):
				m := stats.GetMemStats()

				ch <- &model.MemStats{
					Alloc:      m.Alloc,
					TotalAlloc: m.TotalAlloc,
					Sys:        m.Sys,
				}
			}
		}
	}()

	return ch, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
