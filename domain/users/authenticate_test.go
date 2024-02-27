package users

import (
	"context"
	"errors"
	"testing"
)

func TestUserNotFound(t *testing.T) {
	service := NewUserService(nil, nil, &MockUserRepository{})

	input := AuthenticateUser{
		Email:    "some@email.com",
		Password: "some-password",
	}

	_, err := service.AuthenticateUser(input)
	if err == nil {
		t.Errorf("Expected error to be present but was nil.")
	}

	var userErr *AuthenticateUserError
	if !errors.As(err, &userErr) {
		t.Errorf("Expected error to be 'AuthenticateUserError' but it was '%s'.", err)
	}

	if userErr.Code != EmailAndPasswordDoesNotMatch {
		t.Errorf("Expected error code to be '%s' but it was '%s'.", EmailAndPasswordDoesNotMatch, userErr.Code)
	}
}

func TestIncorrectPassword(t *testing.T) {
	service := NewUserService(nil, nil, &MockUserRepository{})

	input := AuthenticateUser{
		Email:    "existing@email.com",
		Password: "",
	}

	_, err := service.AuthenticateUser(input)
	if err != nil {
		var userErr *AuthenticateUserError
		if errors.As(err, &userErr) && userErr.Code == EmailAndPasswordDoesNotMatch {
			return
		}
	}

	t.Errorf("Expected error to be '%s' but it was '%s'.", EmailAndPasswordDoesNotMatch, err)
}

func TestCorrectPassword(t *testing.T) {
	service := NewUserService(nil, nil, &MockUserRepository{})

	input := AuthenticateUser{
		Email:    "existing@email.com",
		Password: "123456",
	}

	user, err := service.AuthenticateUser(input)
	if user != nil && err == nil {
		return
	}

	t.Errorf("Expected error to be nil but it was '%s.", err)
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
