package repository

import (
	"github.com/golangTroshin/gophermat/internal/model"
	"gorm.io/gorm"
)

type UserBalanceRepository struct {
	db *gorm.DB
}

func NewUserBalanceRepository(db *gorm.DB) *UserBalanceRepository {
	return &UserBalanceRepository{db}
}

func (r *UserBalanceRepository) CreateUserBalance(userBalance *model.UserBalance) error {
	return r.db.Create(userBalance).Error
}

func (r *UserBalanceRepository) UpdateUserBalance(userBalance *model.UserBalance) error {
	return r.db.Save(userBalance).Error
}
