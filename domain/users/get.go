package users

import (
	"context"
)

func (service *UserService) GetUser(id string) (*User, error) {
	return service.userRepository.GetUserByID(context.Background(), id)
}
