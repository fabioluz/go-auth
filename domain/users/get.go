package users

import (
	"context"
)

func (service *UserService) Get(id string) (*User, error) {
	return service.userRepository.GetByID(context.Background(), id)
}
