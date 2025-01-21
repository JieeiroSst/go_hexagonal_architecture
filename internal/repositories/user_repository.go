package repositories

import (
	"time"

	"github.com/JIeeiroSst/hex/internal/core/domain"
	"github.com/JIeeiroSst/hex/internal/core/ports"
	"github.com/JIeeiroSst/hex/internal/infrastructure/database"
)

type userRepository struct {
	db    *database.PostgresDB
	cache ports.CacheRepository
}

type UserModel struct {
	ID           string `gorm:"primaryKey"`
	Name         string
	Email        string `gorm:"uniqueIndex"`
	Password     string
	LastActiveAt time.Time
	CreatedAt    time.Time
}

func NewUserRepository(db *database.PostgresDB, cache ports.CacheRepository) ports.UserRepository {
	return &userRepository{
		db:    db,
		cache: cache,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	model := &UserModel{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		LastActiveAt: user.LastActiveAt,
		CreatedAt:    user.CreatedAt,
	}
	return r.db.Master.Create(model).Error
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	var model UserModel
	if err := r.db.Slave.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           model.ID,
		Name:         model.Name,
		Email:        model.Email,
		Password:     model.Password,
		LastActiveAt: model.LastActiveAt,
		CreatedAt:    model.CreatedAt,
	}, nil
}

func (r *userRepository) Update(user *domain.User) error {
	model := &UserModel{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		LastActiveAt: user.LastActiveAt,
	}
	return r.db.Master.Save(model).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Master.Delete(&UserModel{}, "id = ?", id).Error
}

func (r *userRepository) GetInactiveUsers(threshold time.Time) ([]*domain.User, error) {
	var models []UserModel
	if err := r.db.Slave.Where("last_active_at < ?", threshold).Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = &domain.User{
			ID:           model.ID,
			Name:         model.Name,
			Email:        model.Email,
			Password:     model.Password,
			LastActiveAt: model.LastActiveAt,
			CreatedAt:    model.CreatedAt,
		}
	}
	return users, nil
}
