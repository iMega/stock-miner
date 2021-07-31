package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	sqlTool "github.com/imega/stock-miner/sql"
)

func (s *Storage) GetUser(ctx context.Context) (domain.User, error) {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return domain.User{}, contexkey.ErrExtractEmail
	}

	q := `select name, avatar, role from user where email = ? and delete=0`
	row := s.db.QueryRowContext(ctx, q, email)

	var name, avatar, role string
	if err := row.Scan(&name, &avatar, &role); err != nil {
		return domain.User{}, fmt.Errorf("failed getting user, %w", err)
	}

	return domain.User{
		Email:  email,
		Name:   name,
		Avatar: avatar,
		Role:   role,
	}, nil
}

func (s *Storage) CreateUser(ctx context.Context, user domain.User) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return contexkey.ErrExtractEmail
	}

	createUserTx := func(ctx context.Context, tx *sql.Tx) error {
		q := `insert into user (email, name, avatar, id, role, create_at)
                values(?,?,?,?,?,?)`

		_, err := tx.ExecContext(
			ctx,
			q,
			email,
			user.Name,
			user.Avatar,
			user.ID,
			user.Role,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to create user, %w", err)
		}

		return nil
	}

	wrapper := sqlTool.TxWrapper{DB: s.db}
	if err := wrapper.Transaction(ctx, nil, createUserTx); err != nil {
		return fmt.Errorf("failed to execute transaction, %w", err)
	}

	return nil
}

func (s *Storage) RemoveUser(ctx context.Context, user domain.User) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return contexkey.ErrExtractEmail
	}

	q := `update user set delete=1 where email = ?`
	if _, err := s.db.ExecContext(ctx, q, email); err != nil {
		return fmt.Errorf("failed to remove user, %w", err)
	}

	return nil
}

func userTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS user (
        email VARCHAR(64) PRIMARY KEY,
        name VARCHAR(64),
        avatar VARCHAR(200),
        id VARCHAR(200),
        deleted INTEGER DEFAULT 0,
        role CHAR(4),
        create_at DATETIME NOT NULL
    )`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to execute query, %w", err)
	}

	return nil
}
