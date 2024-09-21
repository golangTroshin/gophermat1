package service

import (
	"errors"

	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/repository"
)

type BalanceService interface {
	GetUserWithBalance(userID uint) (*model.User, error)
	GetUserBalance(userID uint) (*model.UserBalance, error)
	UpdateUserBalance(user *model.User) error
}

type balanceService struct {
	userRepo    repository.UserRepository
	balanceRepo repository.UserBalanceRepository
}

func NewBalanceService(userRepo repository.UserRepository, balanceRepo repository.UserBalanceRepository) BalanceService {
	return &balanceService{userRepo, balanceRepo}
}

func (s *balanceService) GetUserWithBalance(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetUserWithBalanceByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *balanceService) GetUserBalance(userID uint) (*model.UserBalance, error) {
	user, err := s.userRepo.GetUserWithBalanceByID(userID)
	if err != nil {
		return nil, err
	}

	if user.BalanceID == 0 {
		return nil, errors.New("user balance not found")
	}
	return &user.Balance, nil
}

func (s *balanceService) UpdateUserBalance(user *model.User) error {
	return s.balanceRepo.UpdateUserBalance(&user.Balance)
}
