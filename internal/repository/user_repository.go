package repository

import (
	"github.com/golangTroshin/gophermat/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByLogin(login string) (*model.User, error)
	GetUserWithBalanceByID(userID uint) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByLogin(login string) (*model.User, error) {
	var user model.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserWithBalanceByID(userID uint) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Balance").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
