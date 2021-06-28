package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	tools "github.com/imega/stock-miner/sql"
)

func (s *Storage) Slot(ctx context.Context, figi string) ([]domain.Slot, error) {
	var result []domain.Slot

	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return result, contexkey.ErrExtractEmail
	}

	q := `select
            slot_id,
            id,
            ticker,
            figi,
            start_price,
            change_price,
            buying_price,
            target_price,
            profit,
            qty,
            amount_spent,
            target_amount,
            total_profit,
            currency
        from slot
        where email = ?
    `

	if figi != "" {
		q += "and figi = ?"
	}

	rows, err := s.db.QueryContext(ctx, q, email, figi)
	if err != nil {
		return result, fmt.Errorf("failed to execute query, %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		slot := domain.Slot{
			Email: email,
		}

		err := rows.Scan(
			&slot.SlotID,
			&slot.ID,
			&slot.StockItem.Ticker,
			&slot.StockItem.FIGI,
			&slot.StartPrice,
			&slot.ChangePrice,
			&slot.BuyingPrice,
			&slot.TargetPrice,
			&slot.Profit,
			&slot.Qty,
			&slot.AmountSpent,
			&slot.TargetAmount,
			&slot.TotalProfit,
			&slot.StockItem.Currency,
		)
		if err != nil {
			return result, fmt.Errorf("failed to execute query, %w", err)
		}

		result = append(result, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed getting row, %w", err)
	}

	return result, nil
}

func (s *Storage) addSlot(ctx context.Context, t domain.Slot) error {
	q := `insert into slot (
            email,
            slot_id,
            id,
            ticker,
            figi,
            start_price,
            change_price,
            buying_price,
            target_price,
            profit,
            qty,
            amount_spent,
            target_amount,
            total_profit,
            currency
        ) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := s.db.ExecContext(
		ctx,
		q,
		t.Email,
		t.SlotID,
		t.ID,
		t.StockItem.Ticker,
		t.StockItem.FIGI,
		t.StartPrice,
		t.ChangePrice,
		t.BuyingPrice,
		t.TargetPrice,
		t.Profit,
		t.Qty,
		t.AmountSpent,
		t.TargetPrice,
		t.TotalProfit,
		t.StockItem.Currency,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query, %w", err)
	}

	return nil
}

func (s *Storage) BuyStockItem(ctx context.Context, tr domain.Transaction) error {
	tx := func(ctx context.Context, tx *sql.Tx) error {
		if err := s.addSlot(ctx, tr.Slot); err != nil {
			return err
		}

		return s.buyTransaction(ctx, tr)
	}

	wrapper := tools.TxWrapper{s.db}
	if err := wrapper.Transaction(ctx, nil, tx); err != nil {
		return fmt.Errorf("failed to execute transaction, %w", err)
	}

	return nil
}

func (s *Storage) updateSlot(ctx context.Context, t domain.Slot) error {
	q := `
        update slot
        set slot_id = ?,
            ticker = ?,
            figi = ?,
            start_price = ?,
            change_price = ?,
            buying_price = ?,
            target_price = ?,
            profit = ?,
            qty = ?,
            amount_spent = ?,
            target_amount = ?,
            total_profit = ?
        where email = ?
          and id = ?
    `

	_, err := s.db.ExecContext(
		ctx,
		q,
		t.SlotID,
		t.StockItem.Ticker,
		t.StockItem.FIGI,
		t.StartPrice,
		t.ChangePrice,
		t.BuyingPrice,
		t.TargetPrice,
		t.Profit,
		t.Qty,
		t.AmountSpent,
		t.TargetPrice,
		t.TotalProfit,
		t.Email,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update slot, %w", err)
	}

	return nil
}

func (s *Storage) deleteSlot(ctx context.Context, t domain.Slot) error {
	q := `
        delete from slot
        where email = ?
          and id = ?
    `

	_, err := s.db.ExecContext(ctx, q, t.Email, t.ID)
	if err != nil {
		return fmt.Errorf("failed to delete slot, %w", err)
	}

	return nil
}

func slotTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS slot (
        email VARCHAR(64) NOT NULL,
        slot_id INT NOT NULL,
        id VARCHAR(64) NOT NULL,
        ticker VARCHAR(64) NOT NULL,
        figi VARCHAR(200) NOT NULL,

        start_price FLOAT NOT NULL,
        change_price FLOAT NOT NULL,
        buying_price FLOAT NOT NULL,
        target_price FLOAT NOT NULL,
        profit FLOAT NOT NULL,

        qty INT,
        amount_spent FLOAT NOT NULL,
        target_amount FLOAT NOT NULL,
        total_profit FLOAT NOT NULL,

        currency VARCHAR(64) NOT NULL,

        CONSTRAINT pair PRIMARY KEY (email, ticker, id)
    )`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to execute query, %w", err)
	}

	return nil
}

func slotTableMigrate(ctx context.Context, tx *sql.Tx, ti tableInfo) error {
	return nil
}
