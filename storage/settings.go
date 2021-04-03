package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
)

func (s *Storage) Settings(ctx context.Context) (domain.Settings, error) {
	var result domain.Settings

	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return result, fmt.Errorf("failed to extract user from context")
	}

	q := "select doc from settings where email = ?"

	var doc string
	if err := s.db.QueryRowContext(ctx, q, email).Scan(&doc); err != nil {
		if err != sql.ErrNoRows {
			return result, fmt.Errorf("failed to scan settings, %s", err)
		}
		return result, nil
	}

	if err := json.Unmarshal([]byte(doc), &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal settings, %s", err)
	}

	return result, nil
}

func (s *Storage) SaveSettings(ctx context.Context, set domain.Settings) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	b, err := json.Marshal(set)
	if err != nil {
		return fmt.Errorf("failed to marshal settings, %s", err)
	}

	q := "insert into settings (email, doc) values (?,?) on conflict(email) do update set doc = ?"
	if _, err := s.db.ExecContext(ctx, q, email, string(b), string(b)); err != nil {
		return fmt.Errorf("failed to save settings, %s", err)
	}

	return nil
}

func settingsTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS settings (
        email VARCHAR(64) PRIMARY KEY,
        doc TEXT
    )`

	_, err := tx.ExecContext(ctx, q)

	return err
}
