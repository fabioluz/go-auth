package users

import "context"

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
	Name           string `json:"name"`
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	InsertUser(ctx context.Context, input ValidCreateUser) (*User, error)
	UpdateUser(ctx context.Context, id string, input ValidUpdateUser) error
}
