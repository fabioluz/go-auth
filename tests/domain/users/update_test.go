package tests

import (
	"auth/domain/users"
	"errors"
	"testing"
)

func TestInvalidUpdateUserInput(t *testing.T) {
	var tests = []struct {
		name  string
		input users.UpdateUser
		want  users.UpdateUserError
	}{
		{
			name: "Empty name",
			input: users.UpdateUser{
				Name: "",
			},
			want: users.UpdateUserError{
				Code: users.NameIsEmpty,
			},
		},
	}

	service := users.NewUserService(&MockTransactionRepository{}, &MockLogRepository{}, &MockUserRepository{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateUser("", tt.input)
			if err != nil {
				var userErr *users.UpdateUserError
				if errors.As(err, &userErr) && userErr.Code == tt.want.Code {
					return
				}
			}

			t.Errorf("got %s, want nil", err)
		})
	}
}

func TestValidUpdateUserInput(t *testing.T) {
	service := users.NewUserService(&MockTransactionRepository{}, &MockLogRepository{}, &MockUserRepository{})
	validInput := users.UpdateUser{
		Name: "Test",
	}

	err := service.UpdateUser("", validInput)
	if err != nil {
		t.Errorf("Update user returned error: %s", err)
		return
	}
}
