package ports

import "github.com/JIeeiroSst/hex/internal/core/domain"

type UserService interface {
	CreateUser(user *domain.User) error
	GetUser(id string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id string) error
}
