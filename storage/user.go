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

	q := `select name, avatar, role, id from user where email = ? and deleted=0`
	row := s.db.QueryRowContext(ctx, q, email)

	u := domain.User{Email: email}
	if err := row.Scan(&u.Name, &u.Avatar, &u.Role, &u.ID); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to scan record, %w", err)
	}

	return u, nil
}

func (s *Storage) CreateUser(ctx context.Context, user domain.User) error {
	createUserTx := func(ctx context.Context, tx *sql.Tx) error {
		q := `insert into user (email, name, avatar, id, role, create_at)
                values(?,?,?,?,?,?)`

		_, err := tx.ExecContext(
			ctx,
			q,
			user.Email,
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
	q := `update user set deleted=1 where email = ?`
	if _, err := s.db.ExecContext(ctx, q, user.Email); err != nil {
		return fmt.Errorf("failed to remove user, %w", err)
	}

	return nil
}

func (s *Storage) UpdateUser(ctx context.Context, user domain.User) error {
	q := `update user set name=?, avatar=? where email = ?`
	_, err := s.db.ExecContext(ctx, q, user.Name, user.Avatar, user.Email)
	if err != nil {
		return fmt.Errorf("failed to remove user, %w", err)
	}

	return nil
}

func (s *Storage) Users(ctx context.Context) ([]domain.User, error) {
	q := `select id, name, avatar, role, email from user where deleted=0`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed getting users, %w", err)
	}
	defer rows.Close()

	result := []domain.User{}
	for rows.Next() {
		var u domain.User

		err := rows.Scan(&u.ID, &u.Name, &u.Avatar, &u.Role, &u.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to scan, %w", err)
		}

		result = append(result, u)
	}

	return result, nil
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
