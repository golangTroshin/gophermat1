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

type AuthService interface {
	RegisterUser(login, password string) (*model.User, error)
	AuthenticateUser(login, password string) (*model.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	balanceRepo repository.UserBalanceRepository
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserNotFount = errors.New("user not found")
)

const tokenExp = time.Hour * 3
const SecretKey = "supersecretkey"
const CookieAuthToken = "auth_token"

func NewAuthService(userRepo repository.UserRepository, balanceRepo repository.UserBalanceRepository) AuthService {
	return &authService{userRepo, balanceRepo}
}

func (s *authService) RegisterUser(login, password string) (*model.User, error) {
	_, err := s.userRepo.GetUserByLogin(login)
	if err == nil {
		return &model.User{}, ErrUserExists
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

func (s *authService) AuthenticateUser(login, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return user, ErrUserNotFount
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return user, ErrInvalidCreds
	}

	return user, nil
}

func SetAuthCookie(userID uint, w http.ResponseWriter) {
	token, err := buildJWTString(userID)

	if err != nil {
		log.Printf("buildJWTString error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: CookieAuthToken, Value: token})
}

func buildJWTString(userID uint) (string, error) {
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
