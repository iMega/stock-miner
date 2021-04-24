package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
)

func (s *Storage) buyTransaction(ctx context.Context, t domain.Transaction) error {
	q := `insert into dealings (
            email,
            id,
            ticker,
            figi,
            start_price,
            change_price,
            buy_order_id,
            buy_at
        ) values (?,?,?,?,?,?,?,?)`

	_, err := s.db.ExecContext(
		ctx,
		q,
		t.Email,
		t.ID,
		t.Ticker,
		t.FIGI,
		t.StartPrice,
		t.ChangePrice,
		t.BuyOrderID,
		t.BuyAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ConfirmBuyTransaction(ctx context.Context, t domain.Transaction) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	q := `update dealings
        set buying_price = ?,
            target_price = ?,
            profit = ?,
            qty = ?,
            amount_spent = ?
        where email = ?
          and id = ?
    `

	_, err := s.db.ExecContext(
		ctx,
		q,
		t.Slot.BuyingPrice,
		t.Slot.TargetPrice,
		t.Slot.Profit,
		t.Slot.Qty,
		t.AmountSpent,
		email,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to confirm buy transaction, %s", err)
	}

	return nil
}

func (s *Storage) SellTransaction(ctx context.Context, t domain.Transaction) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	q := `update dealings
        set sale_price = ?,
            amount_income = ?,
            total_profit = ?,
            sell_order_id = ?,
            duration = ?,
            sell_at = ?
        where email = ?
          and id = ?
    `

	_, err := s.db.ExecContext(
		ctx,
		q,
		t.SalePrice,
		t.AmountIncome,
		t.TotalProfit,
		t.SellOrderID,
		t.Duration,
		t.SellAt,
		email,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to save sell transaction, %s", err)
	}

	return nil
}

func (s *Storage) Dealings(ctx context.Context) ([]domain.Transaction, error) {
	var result []domain.Transaction

	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return result, fmt.Errorf("failed to extract user from context")
	}

	q := `select email,
	        id,
	        ticker,
            figi,
	        start_price,
	        change_price,
	        buying_price,
	        target_price,
	        profit,
	        sale_price,
	        qty,
	        amount_spent,
	        amount_income,
	        total_profit,
            buy_order_id,
            sell_order_id,
	        buy_at,
	        duration,
	        sell_at
	    from dealings
	    where email = ?
	    order by buy_at`
	rows, err := s.db.QueryContext(ctx, q, email)
	defer rows.Close()
	if err != nil {
		if err != sql.ErrNoRows {
			return result, err
		}

		return result, nil
	}

	for rows.Next() {
		var (
			buyingPrice  sql.NullFloat64
			targetPrice  sql.NullFloat64
			profit       sql.NullFloat64
			salePrice    sql.NullFloat64
			qty          sql.NullInt64
			amountSpent  sql.NullFloat64
			amountIncome sql.NullFloat64
			totalProfit  sql.NullFloat64
			sellOrderID  sql.NullString
			duration     sql.NullInt64
			sellAt       sql.NullTime
		)
		t := domain.Transaction{}
		err := rows.Scan(
			&t.Slot.Email,
			&t.Slot.ID,
			&t.StockItem.Ticker,
			&t.StockItem.FIGI,
			&t.Slot.StartPrice,
			&t.Slot.ChangePrice,
			&buyingPrice,
			&targetPrice,
			&profit,
			//
			&salePrice,
			&qty,
			&amountSpent,
			//
			&amountIncome,
			&totalProfit,
			//
			&t.BuyOrderID,
			&sellOrderID,
			//
			&t.BuyAt,
			&duration,
			&sellAt,
		)

		if buyingPrice.Valid {
			t.Slot.BuyingPrice = buyingPrice.Float64
		}

		if targetPrice.Valid {
			t.Slot.TargetPrice = targetPrice.Float64
		}

		if profit.Valid {
			t.Slot.Profit = profit.Float64
		}

		if salePrice.Valid {
			t.SalePrice = salePrice.Float64
		}

		if qty.Valid {
			t.Slot.Qty = int(qty.Int64)
		}

		if amountSpent.Valid {
			t.Slot.AmountSpent = amountSpent.Float64
		}

		if amountIncome.Valid {
			t.AmountIncome = amountIncome.Float64
		}

		if totalProfit.Valid {
			t.Slot.TotalProfit = totalProfit.Float64
		}

		if sellOrderID.Valid {
			t.SellOrderID = sellOrderID.String
		}

		if duration.Valid {
			t.Duration = int(duration.Int64)
		}

		t.SellAt = time.Date(0, 0, 0, 0, 0, 0, 0, &time.Location{})
		if sellAt.Valid {
			t.SellAt = sellAt.Time
		}

		if err != nil {
			return result, err
		}

		result = append(result, t)
	}

	return result, nil
}

func dealingsTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS dealings (
        email VARCHAR(64) NOT NULL,
        id VARCHAR(32) NOT NULL,
        ticker VARCHAR(64) NOT NULL,
        figi VARCHAR(200) NOT NULL,

        start_price FLOAT NOT NULL,
        change_price FLOAT NOT NULL,
        buying_price FLOAT,
        target_price FLOAT,
        profit FLOAT,

        sale_price FLOAT,
        qty INT,
        amount_spent FLOAT,

        amount_income FLOAT,
        total_profit FLOAT,

        buy_order_id VARCHAR(64) NOT NULL,
        sell_order_id VARCHAR(64),

        buy_at DATETIME NOT NULL,
        duration INT,
        sell_at DATETIME,

        CONSTRAINT pair PRIMARY KEY (email, id)
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}
