package users

import (
	"auth/domain/logs"
	"context"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// Input struct for creating a user
type CreateUser struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Name            string `json:"name"`
}

type CreateUserErrorCode string

type CreateUserError struct {
	Code    CreateUserErrorCode `json:"code"`
	Message string              `json:"message"`
}

func (err *CreateUserError) Error() string {
	return err.Message
}

const (
	EmailIsEmpty                CreateUserErrorCode = "email_empty"
	EmailIsInvalid              CreateUserErrorCode = "email_invalid"
	EmailAlreadyInUse           CreateUserErrorCode = "email_in_use"
	PasswordIsInvalid           CreateUserErrorCode = "password_invalid"
	ConfirmPasswordDoesNotMatch CreateUserErrorCode = "confirm_password_does_not_match"
)

type ValidCreateUser struct {
	Email          string
	HashedPassword string
	Name           string
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (service *UserService) validateCreateUser(input CreateUser) (*ValidCreateUser, error) {
	input.Email = strings.TrimSpace(input.Email)
	if len(input.Email) == 0 {
		return nil, &CreateUserError{
			Code:    EmailIsEmpty,
			Message: "Email is empty.",
		}
	}

	if !isValidEmail(input.Email) {
		return nil, &CreateUserError{
			Code:    EmailIsInvalid,
			Message: "Email is empty.",
		}
	}

	if len(input.Password) < 6 {
		return nil, &CreateUserError{
			Code:    PasswordIsInvalid,
			Message: "Password is invalid.",
		}
	}

	if input.Password != input.ConfirmPassword {
		return nil, &CreateUserError{
			Code:    ConfirmPasswordDoesNotMatch,
			Message: "Confirm password does not match.",
		}
	}

	user, err := service.userRepository.GetUserByEmail(context.TODO(), input.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		panic(err)
	}

	if user != nil {
		return nil, &CreateUserError{
			Code:    EmailAlreadyInUse,
			Message: "Email is already in use.",
		}
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		panic(err)
	}

	validUser := &ValidCreateUser{
		Email:          input.Email,
		HashedPassword: hashedPassword,
		Name:           input.Name,
	}

	return validUser, nil
}

func (service *UserService) CreateUser(input CreateUser) (*User, error) {
	validUser, err := service.validateCreateUser(input)
	if err != nil {
		return nil, err
	}

	var createdUser *User
	createdUserError := service.transactionRepository.WithTransaction(func(ctx context.Context) error {
		var err error
		createdUser, err = service.userRepository.InsertUser(ctx, *validUser)
		if err != nil {
			return err
		}

		_, err = service.logRepository.InsertLog(ctx, createdUser.ID, logs.LogOperationCreateUser)
		return err
	})

	return createdUser, createdUserError
}
