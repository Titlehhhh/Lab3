package services

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (int, error)
	Logout(sessionID string) error
}
