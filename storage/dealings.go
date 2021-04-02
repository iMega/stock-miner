package storage

import (
	"context"
	"database/sql"
	"fmt"

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
            buying_price,
            target_price,
            profit,
            qty,
            amount_spent,
            total_profit,
            buy_at
        ) values (?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := s.db.ExecContext(
		ctx,
		q,
		t.Email,
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
		t.TotalProfit,
		t.BuyAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// func (s *Storage) SellTransaction(ctx context.Context, t domain.Transaction) error {
// 	email, ok := contexkey.EmailFromContext(ctx)
// 	if !ok {
// 		return fmt.Errorf("failed to extract user from context")
// 	}

// 	return nil
// }

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
		var t domain.Transaction

		err := rows.Scan(
			&t.Email,
			&t.Ticker,
			&t.FIGI,
			&t.StartPrice,
			&t.ChangePrice,
			&t.BuyingPrice,
			&t.TargetPrice,
			&t.Profit,
			&t.SalePrice,
			&t.Qty,
			&t.AmountSpent,
			&t.AmountIncome,
			&t.TotalProfit,
			&t.BuyAt,
			&t.Duration,
			&t.SellAt,
		)
		if err != nil {
			return result, nil
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
        buying_price FLOAT NOT NULL,
        target_price FLOAT NOT NULL,
        profit FLOAT NOT NULL,
        sale_price FLOAT,

        qty INT,
        amount_spent FLOAT,
        amount_income FLOAT,
        total_profit FLOAT,

        buy_at DATETIME,
        duration INT,
        sell_at DATETIME,

        CONSTRAINT pair PRIMARY KEY (email, id)
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}
