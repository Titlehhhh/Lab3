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

	user1, _ := u.GetByUsernameUnsafe(user.Username)
	if user1 != nil {
		return errors.New("user already exists")
	}

	u.users = append(u.users, user)

	return saveUsers(fileName, u.users)
}

func (u *userRepository) RemoveUser(id int) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) GetByUsernameUnsafe(username string) (*models.User, error) {
	for _, user := range u.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (u *userRepository) GetByUsername(username string) (*models.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.GetByUsernameUnsafe(username)
}

func saveUsers(path string, users []models.User) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(users)
}

func loadUsers(path string) ([]models.User, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []models.User
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewUserRepository() UserRepository {

	users, err := loadUsers(fileName)

	if err != nil {
		users = []models.User{}
	}

	return &userRepository{
		users: users,
	}
}
