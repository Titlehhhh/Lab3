package services

import (
	"backApp/models"
	"backApp/repository"
	"errors"
	"strings"
	"time"
)

type authService struct {
	repo repository.UserRepository
}

func validate(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return errors.New("username cannot be empty")
	}

	if strings.TrimSpace(password) == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func (a *authService) Register(username, password string) error {

	errVal := validate(username, password)
	if errVal != nil {
		return errVal
	}

	user, _ := a.repo.GetByUsername(username)
	if user != nil {
		return errors.New("User already exists")
	}

	newUser := models.User{
		Id:               0, //TODO for DB
		Username:         username,
		PasswordHash:     password, //TODO Hash
		RegistrationDate: time.Now().String(),
	}
	return a.repo.AddUser(newUser)
}

func (a *authService) Login(username, password string) (int, error) {
	errVal := validate(username, password)
	if errVal != nil {
		return 0, errVal
	}

	user, err := a.repo.GetByUsername(username)
	if err != nil {
		return 0, errors.New("User not found")
	}

	// TODO Hash
	if password != user.PasswordHash {
		return 0, errors.New("Password does not match")
	}

	return user.Id, nil
}

func (a *authService) Logout(sessionID string) error {
	return nil
}

func NewAuthService(repo *repository.UserRepository) AuthService {
	return &authService{
		repo: *repo,
	}
}
