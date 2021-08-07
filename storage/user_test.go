package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/imega/stock-miner/contexkey"
	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetUser(t *testing.T) {
	db, close, err := helpers.CreateDB(
		func(ctx context.Context, tx *sql.Tx) error {
			return userTable(ctx, tx)
		},
	)
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := &Storage{
		db: db,
	}

	expected := domain.User{
		Email:  "test@example.com",
		Name:   "name",
		ID:     "id",
		Avatar: "avatar",
		Role:   "role",
	}

	ctx := contexkey.WithEmail(context.Background(), expected.Email)

	if err := s.CreateUser(ctx, expected); err != nil {
		t.Fatalf("failed to create user, %s", err)
	}

	actual, err := s.GetUser(ctx)
	if err != nil {
		t.Fatalf("failed getting user, %s", err)
	}

	assert.Equal(t, expected, actual)
}

func TestStorage_RemoveUser(t *testing.T) {
	db, close, err := helpers.CreateDB(
		func(ctx context.Context, tx *sql.Tx) error {
			return userTable(ctx, tx)
		},
	)
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := &Storage{
		db: db,
	}

	expected := domain.User{
		Email:  "test@example.com",
		Name:   "name",
		ID:     "id",
		Avatar: "avatar",
		Role:   "role",
	}

	ctx := contexkey.WithEmail(context.Background(), expected.Email)

	if err := s.CreateUser(ctx, expected); err != nil {
		t.Fatalf("failed to create user, %s", err)
	}

	if err := s.RemoveUser(ctx, expected); err != nil {
		t.Fatalf("failed to remove user, %s", err)
	}

	_, err = s.GetUser(ctx)
	if err != domain.ErrUserNotFound {
		t.Fatalf("other error, %s", err)
	}
}

func TestStorage_Users(t *testing.T) {
	db, close, err := helpers.CreateDB(
		func(ctx context.Context, tx *sql.Tx) error {
			return userTable(ctx, tx)
		},
	)
	if err != nil {
		t.Fatalf("failed to create database, %s", err)
	}
	defer close()

	s := &Storage{
		db: db,
	}

	expected := domain.User{
		Email:  "test@example.com",
		Name:   "name",
		ID:     "id",
		Avatar: "avatar",
		Role:   "role",
	}

	ctx := contexkey.WithEmail(context.Background(), expected.Email)

	if err := s.CreateUser(ctx, expected); err != nil {
		t.Fatalf("failed to create user, %s", err)
	}

	actual, err := s.Users(ctx)
	if err != nil {
		t.Fatalf("failed getting users, %s", err)
	}

	assert.Equal(t, expected, actual[0])
}
