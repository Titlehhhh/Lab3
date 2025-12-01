package models

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type User struct {
	Id               int
	Name             string
	PasswordHash     string
	RegistrationDate string
}
