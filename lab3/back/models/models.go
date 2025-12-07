package models

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type User struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	PasswordHash     string `json:"password_hash"`
	RegistrationDate string `json:"registration_date"`
}

type UserRegistration struct {
	Username string
	Password string
}

type Session struct{}
