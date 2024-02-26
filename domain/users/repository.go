package users

import "context"

type User struct {
	ID             string
	Email          string
	HashedPassword string
	Name           string
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	InsertUser(ctx context.Context, input ValidCreateUser) (*User, error)
}
