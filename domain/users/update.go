package users

import (
	"auth/domain/logs"
	"context"
)

// Unvalidated input for updating user
type UpdateUser struct {
	Name string `json:"name"`
}

type UpdateUserErrorCode string

type UpdateUserError struct {
	Code    UpdateUserErrorCode `json:"code"`
	Message string              `json:"message"`
}

func (err *UpdateUserError) Error() string {
	return err.Message
}

const (
	NameIsEmpty UpdateUserErrorCode = "name_empty"
)

// Validated input for udpating user
type ValidUpdateUser struct {
	Name string
}

func (service *UserService) validateUpdateUser(input UpdateUser) (*ValidUpdateUser, error) {
	if input.Name == "" {
		return nil, &UpdateUserError{
			Code:    NameIsEmpty,
			Message: "Name is empty.",
		}
	}

	validUser := ValidUpdateUser{
		Name: input.Name,
	}

	return &validUser, nil
}

func (service *UserService) UpdateUser(id string, input UpdateUser) error {
	validUser, err := service.validateUpdateUser(input)
	if err != nil {
		return err
	}

	updateUserError := service.transactionRepository.WithTransaction(func(ctx context.Context) error {
		err := service.userRepository.UpdateUser(ctx, id, *validUser)
		if err != nil {
			return err
		}

		_, err = service.logRepository.InsertLog(ctx, id, logs.LogOperationUpdateUser)
		return err
	})

	return updateUserError
}
