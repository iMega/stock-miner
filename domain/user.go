package domain

import (
	"context"
	"errors"
)

type User struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	ID     string `json:"id"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

var ErrUserNotFound = errors.New("user not found")

type UserStorage interface {
	GetUser(context.Context) (User, error)
	Users(context.Context) ([]User, error)

	CreateUser(context.Context, User) error
	UpdateUser(context.Context, User) error
	RemoveUser(context.Context, User) error
}
