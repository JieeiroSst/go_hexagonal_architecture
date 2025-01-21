package ports

import "github.com/JIeeiroSst/hex/internal/core/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
}
