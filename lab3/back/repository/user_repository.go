package repository

import (
	"backApp/models"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type userRepository struct {
	users []models.User
	mu    sync.Mutex
}

const fileName = "users.json"

func (u *userRepository) AddUser(user models.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users = append(u.users, user)
	data, err := json.MarshalIndent(u.users, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func (u *userRepository) RemoveUser(id int) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) GetByUsername(username string) (*models.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, user := range u.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
