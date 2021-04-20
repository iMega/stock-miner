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
		return result, fmt.Errorf("failed to extract user from context")
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
            total_profit
        from slot
        where email = ? and figi = ?
    `

	rows, err := s.db.QueryContext(ctx, q, email, figi)
	defer rows.Close()
	if err != nil {
		return result, err
	}

	for rows.Next() {
		slot := domain.Slot{
			Email: email,
		}
		err := rows.Scan(
			&slot.SlotID,
			&slot.ID,
			&slot.Ticker,
			&slot.FIGI,
			&slot.StartPrice,
			&slot.ChangePrice,
			&slot.BuyingPrice,
			&slot.TargetPrice,
			&slot.Profit,
			&slot.Qty,
			&slot.AmountSpent,
			&slot.TargetAmount,
			&slot.TotalProfit,
		)
		if err != nil {
			return result, err
		}

		result = append(result, slot)
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
            total_profit
        ) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := s.db.ExecContext(
		ctx,
		q,
		t.Email,
		t.SlotID,
		t.ID,
		t.Ticker,
		t.FIGI,
		t.StartPrice,
		t.ChangePrice,
		t.BuyingPrice,
		t.TargetPrice,
		t.Profit,
		t.Qty,
		t.AmountSpent,
		t.TargetPrice,
		t.TotalProfit,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) BuyStockItem(ctx context.Context, slot domain.Slot) error {
	wrapper := tools.TxWrapper{s.db}
	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.addSlot(ctx, slot); err != nil {
			return err
		}

		t := domain.Transaction{Slot: slot}
		if err := s.buyTransaction(ctx, t); err != nil {
			return err
		}

		return nil
	})
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

        CONSTRAINT pair PRIMARY KEY (email, slot_id)
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}

func slotTableMigrate(ctx context.Context, tx *sql.Tx, ti tableInfo) error {
	return nil
}
