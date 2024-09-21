package service

import (
	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/repository"
)

type WithdrawService interface {
	RecordWithdrawal(withdrawal *model.UserWithdrawal) error
	GetUserWithdrawals(userID uint) ([]model.UserWithdrawal, error)
}

type withdrawService struct {
	withdrawRepo repository.WithdrawRepository
}

func NewWithdrawService(withdrawRepo repository.WithdrawRepository) WithdrawService {
	return &withdrawService{withdrawRepo}
}

func (s *withdrawService) RecordWithdrawal(withdrawal *model.UserWithdrawal) error {
	return s.withdrawRepo.RecordWithdrawal(withdrawal)
}

func (s *withdrawService) GetUserWithdrawals(userID uint) ([]model.UserWithdrawal, error) {
	return s.withdrawRepo.GetUserWithdrawals(userID)
}
