package services

import "backApp/models"

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (models.Session, error)
	Logout(sessionID string) error
}
