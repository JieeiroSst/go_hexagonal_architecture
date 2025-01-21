package services

import (
	"github.com/JIeeiroSst/hex/internal/core/domain"
	"github.com/JIeeiroSst/hex/internal/core/ports"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *domain.User) error {
	return s.userRepo.Create(user)
}

func (s *userService) GetUser(id string) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) UpdateUser(user *domain.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}
