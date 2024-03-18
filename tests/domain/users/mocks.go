package tests

import (
	"auth/domain/logs"
	"auth/domain/users"
	"context"
	"time"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	panic("not implemented")
}

func (repo *MockUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	if email == "existing@email.com" {
		user := users.User{
			ID:             "",
			Email:          "existing@email.com",
			HashedPassword: "$2a$10$rENhbYtCaSrx5vqUoK3UaulMQzUQppnG/upOgBvAtmI3rgOapqo66",
			Name:           "",
		}

		return &user, nil
	}
	return nil, nil
}

func (repo *MockUserRepository) Insert(ctx context.Context, input users.ValidCreateUser) (*users.User, error) {
	return &users.User{
		ID:             "65dc8fa57bd9e61fb8817a09",
		Email:          input.Email,
		HashedPassword: "$2a$10$rENhbYtCaSrx5vqUoK3UaulMQzUQppnG/upOgBvAtmI3rgOapqo66",
		Name:           input.Name,
	}, nil
}

func (repo *MockUserRepository) Update(ctx context.Context, id string, input users.ValidUpdateUser) error {
	return nil
}

type MockTransactionRepository struct{}

func (repo *MockTransactionRepository) WithTransaction(fn func(ctx context.Context) error) error {
	return fn(context.TODO())
}

type MockLogRepository struct{}

func (repo *MockLogRepository) Insert(ctx context.Context, userID string, op logs.LogOperation) (*logs.Log, error) {
	return &logs.Log{
		ID:        "65dc8fa57bd9e61fb8817a0a",
		UserID:    "65dc8fa57bd9e61fb8817a09",
		Operation: op,
		Timestamp: time.Now().UTC(),
	}, nil
}

func (repo *MockLogRepository) Get(ctx context.Context, userID string, pageSize int, after string) ([]logs.Log, error) {
	return []logs.Log{}, nil
}
