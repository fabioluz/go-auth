package users

import (
	"context"
)

type AuthenticateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserErrorCode string

type AuthenticateUserError struct {
	Code    AuthenticateUserErrorCode `json:"code"`
	Message string                    `json:"message"`
}

func (err *AuthenticateUserError) Error() string {
	return err.Message
}

const (
	EmailAndPasswordDoesNotMatch AuthenticateUserErrorCode = "email_and_password_does_not_match"
)

func authenticateError() (*User, error) {
	return nil, &AuthenticateUserError{
		Code:    EmailAndPasswordDoesNotMatch,
		Message: "Email and/or Password does not match.",
	}
}

func (service *UserService) AuthenticateUser(input AuthenticateUser) (*User, error) {
	user, err := service.userRepository.GetUserByEmail(context.Background(), input.Email)
	if err != nil || user == nil {
		return authenticateError()
	}

	if !comparePassword(user.HashedPassword, input.Password) {
		return authenticateError()
	}

	return user, nil
}
