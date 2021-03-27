package domain

import "context"

type User struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	ID     string `json:"id"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

type UserStorage interface {
	GetUser(context.Context) (User, error)

	CreateUser(context.Context, User) error
	RemoveUser(context.Context, User) error
}
