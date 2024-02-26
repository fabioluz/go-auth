package users

import (
	"auth/domain"
	"auth/domain/logs"
)

type UserService struct {
	transactionRepository domain.TransactionRepository
	logRepository         logs.LogRepository
	userRepository        UserRepository
}

func NewUserService(
	transactionRepository domain.TransactionRepository,
	logRepository logs.LogRepository,
	userRepository UserRepository,
) *UserService {

	return &UserService{
		transactionRepository,
		logRepository,
		userRepository,
	}
}
