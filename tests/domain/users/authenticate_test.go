package tests

import (
	"auth/domain/users"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestUserNotFound(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	service := users.NewUserService(validate, nil, nil, &MockUserRepository{})

	input := users.AuthenticateUserInput{
		Email:    "some@email.com",
		Password: "some-password",
	}

	_, err := service.Authenticate(input)
	if err == nil {
		t.Errorf("Expected error to be present but was nil.")
	}

	var userErr *users.AuthenticateUserError
	if !errors.As(err, &userErr) {
		t.Errorf("Expected error to be 'AuthenticateUserError' but it was '%s'.", err)
	}

	if userErr.Code != users.EmailAndPasswordDoesNotMatch {
		t.Errorf("Expected error code to be '%s' but it was '%s'.", users.EmailAndPasswordDoesNotMatch, userErr.Code)
	}
}

func TestIncorrectPassword(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	service := users.NewUserService(validate, nil, nil, &MockUserRepository{})

	input := users.AuthenticateUserInput{
		Email:    "existing@email.com",
		Password: "",
	}

	_, err := service.Authenticate(input)
	if err != nil {
		var userErr *users.AuthenticateUserError
		if errors.As(err, &userErr) && userErr.Code == users.EmailAndPasswordDoesNotMatch {
			return
		}
	}

	t.Errorf("Expected error to be '%s' but it was '%s'.", users.EmailAndPasswordDoesNotMatch, err)
}

func TestCorrectPassword(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	service := users.NewUserService(validate, nil, nil, &MockUserRepository{})

	input := users.AuthenticateUserInput{
		Email:    "existing@email.com",
		Password: "123456",
	}

	user, err := service.Authenticate(input)
	if user != nil && err == nil {
		return
	}

	t.Errorf("Expected error to be nil but it was '%s.", err)
}
