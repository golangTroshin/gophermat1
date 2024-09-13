package repository

import (
	"github.com/golangTroshin/gophermat/internal/model"
	"gorm.io/gorm"
)

type WithdrawRepository interface {
	RecordWithdrawal(withdrawal *model.UserWithdrawal) error
	GetUserWithdrawals(userID uint) ([]model.UserWithdrawal, error)
}

type withdrawRepository struct {
	db *gorm.DB
}

func NewWithdrawRepository(db *gorm.DB) WithdrawRepository {
	return &withdrawRepository{db}
}

func (r *withdrawRepository) RecordWithdrawal(withdrawal *model.UserWithdrawal) error {
	return r.db.Create(withdrawal).Error
}

func (r *withdrawRepository) GetUserWithdrawals(userID uint) ([]model.UserWithdrawal, error) {
	var withdrawals []model.UserWithdrawal
	if err := r.db.Where("user_id = ?", userID).Order("processed_at DESC").Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}
