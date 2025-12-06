package repository

import (
	"backApp/models"
	"errors"
	"sync"
)

type userRepository struct {
	users sync.Map //TODO file
}

const fileName = "users.json"

func (u *userRepository) AddUser(user models.User) error {
	if _, exists := u.users.Load(user.Username); exists {
		return errors.New("user already exists")
	}
	u.users.Store(user.Username, user)
	return nil
}

func (u *userRepository) RemoveUser(id int) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) GetByUsername(username string) (*models.User, error) {
	if u, ok := u.users.Load(username); ok {
		user := u.(models.User)
		return &user, nil
	}
	return nil, errors.New("user not found")
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
