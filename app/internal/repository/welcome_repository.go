package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/yamamoto99/go-template/app/internal/entity"
)

type WelcomeRepository interface {
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

type welcomeRepository struct {
	db *gorm.DB
}

func NewWelcomeRepository(db *gorm.DB) WelcomeRepository {
	return &welcomeRepository{db: db}
}

func (r *welcomeRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
