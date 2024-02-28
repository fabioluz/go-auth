package tests

import (
	"auth/domain/users"
	"errors"
	"testing"
)

func TestInvalidInput(t *testing.T) {
	var tests = []struct {
		name  string
		input users.CreateUser
		want  users.CreateUserError
	}{
		{
			name: "Empty email",
			input: users.CreateUser{
				Email: "",
			},
			want: users.CreateUserError{
				Code: users.EmailIsEmpty,
			},
		},
		{
			name: "Invalid email",
			input: users.CreateUser{
				Email: "test",
			},
			want: users.CreateUserError{
				Code: users.EmailIsInvalid,
			},
		},
		{
			name: "Invalid password",
			input: users.CreateUser{
				Email:           "non-existing@email.com",
				Password:        "12345",
				ConfirmPassword: "12345",
			},
			want: users.CreateUserError{
				Code: users.PasswordIsInvalid,
			},
		},
		{
			name: "Unmatched password",
			input: users.CreateUser{
				Email:           "non-existing@email.com",
				Password:        "123456",
				ConfirmPassword: "12345",
			},
			want: users.CreateUserError{
				Code: users.ConfirmPasswordDoesNotMatch,
			},
		},
		{
			name: "Email in use",
			input: users.CreateUser{
				Email:           "existing@email.com",
				Password:        "123456",
				ConfirmPassword: "123456",
				Name:            "Test",
			},
			want: users.CreateUserError{
				Code: users.EmailAlreadyInUse,
			},
		},
	}

	service := users.NewUserService(&MockTransactionRepository{}, &MockLogRepository{}, &MockUserRepository{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.CreateUser(tt.input)
			if err != nil {
				var userErr *users.CreateUserError
				if errors.As(err, &userErr) && userErr.Code == tt.want.Code {
					return
				}
			}

			t.Errorf("got %s, want %s", user, tt.want)
		})
	}
}

func TestValidInput(t *testing.T) {
	service := users.NewUserService(&MockTransactionRepository{}, &MockLogRepository{}, &MockUserRepository{})
	validInput := users.CreateUser{
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
