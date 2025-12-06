package repository

import "backApp/models"

type UserRepository interface {
	AddUser(models.User) error
	RemoveUser(id int) error
	GetByUsername(username string) (*models.User, error)
}
