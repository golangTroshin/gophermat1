package repository

import (
	"github.com/golangTroshin/gophermat/internal/model"
	"gorm.io/gorm"
)

type UserBalanceRepository interface {
	CreateUserBalance(userBalance *model.UserBalance) error
	UpdateUserBalance(userBalance *model.UserBalance) error
}

type userBalanceRepository struct {
	db *gorm.DB
}

func NewUserBalanceRepository(db *gorm.DB) UserBalanceRepository {
	return &userBalanceRepository{db}
}

func (r *userBalanceRepository) CreateUserBalance(userBalance *model.UserBalance) error {
	return r.db.Create(userBalance).Error
}

func (r *userBalanceRepository) UpdateUserBalance(userBalance *model.UserBalance) error {
	return r.db.Save(userBalance).Error
}
