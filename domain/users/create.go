package users

import (
	"auth/domain/logs"
	"context"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Unvalidated Input for creating user
type CreateUser struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"gte=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"eqfield=Password"`
	Name            string `json:"name" validate:"required"`
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

// Validated input for creating user
type ValidCreateUser struct {
	Email          string
	HashedPassword string
	Name           string
}

func (service *UserService) validateCreateUser(input CreateUser) (*ValidCreateUser, error) {
	input.Email = strings.TrimSpace(input.Email)

	err := service.validate.Struct(&input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				switch err.Tag() {
				case "required":
					return nil, &CreateUserError{
						Code:    EmailIsEmpty,
						Message: "Email is empty.",
					}
				case "email":
					return nil, &CreateUserError{
						Code:    EmailIsInvalid,
						Message: "Email is invalid.",
					}
				}
			case "Password":
				switch err.Tag() {
				case "required":
					return nil, &CreateUserError{
						Code:    PasswordIsInvalid,
						Message: "Password is invalid.",
					}
				case "gte":
					return nil, &CreateUserError{
						Code:    PasswordIsInvalid,
						Message: "Password is invalid.",
					}
				}
			case "ConfirmPassword":
				switch err.Tag() {
				case "eqfield":
					return nil, &CreateUserError{
						Code:    ConfirmPasswordDoesNotMatch,
						Message: "Confirm password does not match.",
					}

				}
			case "Name":
				switch err.Tag() {
				case "eqfield":
					return nil, &CreateUserError{
						Code:    "name_is_empty",
						Message: "Name is empty.",
					}

				}
			}
		}
	}

	user, err := service.userRepository.GetUserByEmail(context.TODO(), input.Email)
	if err != nil {
		panic(err)
	}

	if user != nil {
		return nil, &CreateUserError{
			Code:    EmailAlreadyInUse,
			Message: "Email is already in use.",
		}
	}

	hashedPassword, err := hashPassword(input.Password)
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
