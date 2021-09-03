package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/uuid"
)

func (b *Broker) branchBuyOrSell(msg domain.Message) error {
	slot := msg.Transaction.Slot

	frame, err := b.SMAStack.Get(slot.StockItem.Ticker)
	if err != nil {
		return fmt.Errorf("failed getting frame from stack, %w", err)
	}

	h, l := frame.RangeHL()
	if h == 0 || l == 0 {
		return fmt.Errorf("%w, figi=%s", errRangeIsZero, slot.StockItem.FIGI)
	}

	if !frame.IsFull() {
		return errFrameNotFull
	}

	ctx, err := b.contextWithCreds(context.Background(), slot.Email)
	if err != nil {
		return fmt.Errorf("failed getting creds, %w", err)
	}

	settings, err := b.SettingsStorage.Settings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting settings, %w", err)
	}

	slots, err := b.Stack.Slot(ctx, slot.StockItem.FIGI)
	if err != nil {
		return fmt.Errorf("failed getting slot, %w", err)
	}

	sellSlots := getItemsForSale(slots, frame.Last(), frame.Prev())
	for _, slot := range sellSlots {
		b.Traffic.SellCh <- slot
	}

	if settings.Slot.Volume <= len(slots) {
		return nil
	}

	maxPurchasePrice := 0.0

	if slot.StockItem.MaxPrice > 0 {
		maxPurchasePrice = slot.StockItem.MaxPrice
	}

	minPrice := minBuyingPrice(slots, maxPurchasePrice)
	if minPrice == 0 {
		return nil
	}

	if minPrice-settings.Slot.ModificatorMinPrice >= msg.Price {
		return nil
	}

	isTrendUp, err := b.SMAStack.IsTrendUp(slot.StockItem.Ticker)
	if err != nil {
		return fmt.Errorf("failed getting trend, %w", err)
	}

	if isTrendUp {
		return nil
	}

	if !priceInRange(frame, msg.Price) {
		return nil
	}

	// buy
	emptyTr := domain.Transaction{
		Slot: domain.Slot{
			ID:          uuid.NewID().String(),
			Email:       msg.Transaction.Email,
			StockItem:   slot.StockItem,
			SlotID:      len(slots) + 1,
			StartPrice:  frame.Prev(),
			ChangePrice: frame.Last(),
			Qty:         1,
		},
		BuyAt: time.Now(),
	}

	tr, err := b.buy(ctx, emptyTr)
	if err != nil {
		return fmt.Errorf("failed to buy stock item, %w", err)
	}

	b.Traffic.ConfirmBuyCh <- domain.Message{Transaction: tr}

	return nil
}
