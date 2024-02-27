package users

import (
	"auth/domain/logs"
	"context"
	"errors"
	"testing"
	"time"
)

func TestInvalidInput(t *testing.T) {
	var tests = []struct {
		name  string
		input CreateUser
		want  CreateUserError
	}{
		{
			name: "Empty email",
			input: CreateUser{
				Email: "",
			},
			want: CreateUserError{
				Code: EmailIsEmpty,
			},
		},
		{
			name: "Invalid email",
			input: CreateUser{
				Email: "test",
			},
			want: CreateUserError{
				Code: EmailIsInvalid,
			},
		},
		{
			name: "Invalid password",
			input: CreateUser{
				Email:           "non-existing@email.com",
				Password:        "12345",
				ConfirmPassword: "12345",
			},
			want: CreateUserError{
				Code: PasswordIsInvalid,
			},
		},
		{
			name: "Unmatched password",
			input: CreateUser{
				Email:           "non-existing@email.com",
				Password:        "123456",
				ConfirmPassword: "12345",
			},
			want: CreateUserError{
				Code: ConfirmPasswordDoesNotMatch,
			},
		},
		{
			name: "Email in use",
			input: CreateUser{
				Email:           "existing@email.com",
				Password:        "123456",
				ConfirmPassword: "123456",
				Name:            "Test",
			},
			want: CreateUserError{
				Code: EmailAlreadyInUse,
			},
		},
	}

	service := NewUserService(&MockTransactionRepository{}, &MockLogRepository{}, &MockUserRepository{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.CreateUser(tt.input)
			if err != nil {
				var userErr *CreateUserError
				if errors.As(err, &userErr) && userErr.Code == tt.want.Code {
					return
				}
			}

			t.Errorf("got %s, want %s", user, tt.want)
		})
	}
}

func TestValidInput(t *testing.T) {
	service := NewUserService(&MockTransactionRepository{}, &MockLogRepository{}, &MockUserRepository{})
	validInput := CreateUser{
		Email:           "non-existing@email.com",
		Password:        "123456",
		ConfirmPassword: "123456",
		Name:            "Test",
	}

	user, err := service.CreateUser(validInput)
	if err != nil {
		t.Errorf("create user returned error: %s", err)
		return
	}

	if user == nil {
		t.Errorf("create user is nil")
		return
	}

	if user.Email != validInput.Email {
		t.Errorf("got %s, want %s", user.Email, validInput.Email)
	}

	if user.Name != validInput.Name {
		t.Errorf("got %s, want %s", user.Name, validInput.Name)
	}
}

type MockTransactionRepository struct{}

func (repo *MockTransactionRepository) WithTransaction(fn func(ctx context.Context) error) error {
	return fn(context.TODO())
}

type MockLogRepository struct{}

func (repo *MockLogRepository) InsertLog(ctx context.Context, userID string, op logs.LogOperation) (*logs.Log, error) {
	return &logs.Log{
		ID:        "65dc8fa57bd9e61fb8817a0a",
		UserID:    "65dc8fa57bd9e61fb8817a09",
		Operation: op,
		Timestamp: time.Now().UTC(),
	}, nil
}

func (repo *MockLogRepository) GetLogs(ctx context.Context, userID string, pageSize int, after string) ([]logs.Log, error) {
	return []logs.Log{}, nil
}
