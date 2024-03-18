package users

import (
	"auth/domain"
	"auth/domain/logs"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	validate              *validator.Validate
	transactionRepository domain.TransactionRepository
	logRepository         logs.LogRepository
	userRepository        UserRepository
}

func NewUserService(
	validate *validator.Validate,
	transactionRepository domain.TransactionRepository,
	logRepository logs.LogRepository,
	userRepository UserRepository,
) *UserService {

	return &UserService{
		validate,
		transactionRepository,
		logRepository,
		userRepository,
	}
}
