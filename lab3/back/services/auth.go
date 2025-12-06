package services

import (
	"backApp/models"
	"backApp/repository"
	"errors"
	"time"
)

type authService struct {
	repo repository.UserRepository
}

func (a *authService) Register(username, password string) error {
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

func (a *authService) Login(username, password string) (models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (a *authService) Logout(sessionID string) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthService(repo *repository.UserRepository) AuthService {
	return &authService{
		repo: *repo,
	}
}
