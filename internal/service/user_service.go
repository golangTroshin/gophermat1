package service

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/repository"
	"github.com/golangTroshin/gophermat/internal/utils"
)

type UserService struct {
	userRepo     *repository.UserRepository
	balanceRepo  *repository.UserBalanceRepository
	withdrawRepo *repository.WithdrawRepository
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint
}

const tokenExp = time.Hour * 3
const SecretKey = "supersecretkey"
const CookieAuthToken = "auth_token"

func NewUserService(userRepo *repository.UserRepository, balanceRepo *repository.UserBalanceRepository, withdrawRepo *repository.WithdrawRepository) *UserService {
	return &UserService{userRepo, balanceRepo, withdrawRepo}
}

func (s *UserService) RegisterUser(login, password string) (*model.User, error) {
	_, err := s.userRepo.GetUserByLogin(login)
	if err == nil {
		return &model.User{}, errors.New("login already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return &model.User{}, err
	}

	balance := &model.UserBalance{}
	s.balanceRepo.CreateUserBalance(balance)

	user := &model.User{
		Login:     login,
		Password:  hashedPassword,
		BalanceID: balance.ID,
	}
	result := s.userRepo.CreateUser(user)

	return user, result
}

func (s *UserService) AuthenticateUser(login, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return user, errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return user, errors.New("invalid password")
	}

	return user, nil
}

func BuildJWTString(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetAuthCookie(userID uint, w http.ResponseWriter) {
	token, err := BuildJWTString(userID)

	if err != nil {
		log.Printf("BuildJWTString error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: CookieAuthToken, Value: token})
}

func (s *UserService) GetUserWithBalanceBalance(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetUserWithBalanceByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserBalance(userID uint) (*model.UserBalance, error) {
	user, err := s.userRepo.GetUserWithBalanceByID(userID)
	if err != nil {
		return nil, err
	}

	if user.BalanceID == 0 {
		return nil, errors.New("user balance not found")
	}
	return &user.Balance, nil
}

func (s *UserService) UpdateUserBalance(user *model.User) error {
	return s.balanceRepo.UpdateUserBalance(&user.Balance)
}

func (s *UserService) RecordWithdrawal(withdrawal *model.UserWithdrawal) error {
	return s.withdrawRepo.RecordWithdrawal(withdrawal)
}

func (s *UserService) GetUserWithdrawals(userID uint) ([]model.UserWithdrawal, error) {
	return s.withdrawRepo.GetUserWithdrawals(userID)
}
