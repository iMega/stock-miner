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
		return domain.User{}, fmt.Errorf("failed to extract user from context")
	}

	q := `select name, avatar, role from user`
	row := s.db.QueryRowContext(ctx, q, email)
	var name, avatar, role string
	if err := row.Scan(&name, &avatar, &role); err != nil {
		return domain.User{}, fmt.Errorf("failed getting user")
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
		return fmt.Errorf("failed to extract user from context")
	}
	wrapper := sqlTool.TxWrapper{s.db}

	return wrapper.Transaction(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		userQuery := `insert into user (email, name, avatar, id, role, create_at) values(?,?,?,?,?,?)`
		_, err := tx.ExecContext(ctx, userQuery, email, user.Name, user.Avatar, user.ID, user.Role, time.Now())
		if err != nil {
			return fmt.Errorf("failed to create user, %s", err)
		}

		// settingsQuery := ``
		// if _, err := tx.ExecContext(ctx, settingsQuery, email, "{}"); err != nil {
		// 	return fmt.Errorf("failed to create user, %s", err)
		// }

		return nil
	})
}

func (s *Storage) RemoveUser(ctx context.Context, user domain.User) error {
	email, ok := contexkey.EmailFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract user from context")
	}

	q := `update user set delete=1 where email = ?`
	if _, err := s.db.ExecContext(ctx, q, email); err != nil {
		return fmt.Errorf("failed to remove user, %s", err)
	}

	return nil
}
