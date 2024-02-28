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

type MockUserRepository struct{}

func (repo *MockUserRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	panic("not implemented")
}

func (repo *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "existing@email.com" {
		user := User{
			ID:             "",
			Email:          "existing@email.com",
			HashedPassword: "$2a$10$rENhbYtCaSrx5vqUoK3UaulMQzUQppnG/upOgBvAtmI3rgOapqo66",
			Name:           "",
		}

		return &user, nil
	}
	return nil, nil
}

func (repo *MockUserRepository) InsertUser(ctx context.Context, input ValidCreateUser) (*User, error) {
	return &User{
		ID:             "65dc8fa57bd9e61fb8817a09",
		Email:          input.Email,
		HashedPassword: "$2a$10$rENhbYtCaSrx5vqUoK3UaulMQzUQppnG/upOgBvAtmI3rgOapqo66",
		Name:           input.Name,
	}, nil
}

func (repo *MockUserRepository) UpdateUser(ctx context.Context, id string, input ValidUpdateUser) error {
	panic("not implemented")
}
