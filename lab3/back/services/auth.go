package services

import "backApp/models"

type authService struct {
}

func (a *authService) Register(username, password string) error {
	//TODO implement me
	panic("implement me")
}

func (a *authService) Login(username, password string) (models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (a *authService) Logout(sessionID string) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthService() AuthService {
	return &authService{}
}
