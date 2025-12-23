package repository

import (
	"backApp/models"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

type userRepository struct {
	users []models.User
	mu    sync.Mutex
}

const fileName = "users.json"

func (u *userRepository) AddUser(user models.User) error {

	//Validate

	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username cannot be empty")
	}

	if strings.TrimSpace(user.PasswordHash) == "" {
		return errors.New("password cannot be empty")
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	user1, _ := u.GetByUsernameUnsafe(user.Username)
	if user1 != nil {
		return errors.New("user already exists")
	}

	if user.Id == 0 {
		user.Id = u.nextID()
	}

	if strings.TrimSpace(user.Role) == "" {
		if len(u.users) == 0 {
			user.Role = "admin"
		} else {
			user.Role = "user"
		}
	}

	u.users = append(u.users, user)

	return saveUsers(fileName, u.users)
}

func (u *userRepository) RemoveUser(id int) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	for index, user := range u.users {
		if user.Id == id {
			u.users = append(u.users[:index], u.users[index+1:]...)
			return saveUsers(fileName, u.users)
		}
	}
	return errors.New("user not found")
}

func (u *userRepository) GetByUsernameUnsafe(username string) (*models.User, error) {

	for index, user := range u.users {
		if user.Username == username {
			return &u.users[index], nil
		}
	}
	return nil, errors.New("user not found")
}

func (u *userRepository) GetByUsername(username string) (*models.User, error) {

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	return u.GetByUsernameUnsafe(username)
}

func (u *userRepository) GetByID(id int) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("id cannot be empty")
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	for index, user := range u.users {
		if user.Id == id {
			return &u.users[index], nil
		}
	}
	return nil, errors.New("user not found")
}

func (u *userRepository) nextID() int {
	maxID := 0
	for _, user := range u.users {
		if user.Id > maxID {
			maxID = user.Id
		}
	}
	return maxID + 1
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
