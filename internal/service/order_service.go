package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/golangTroshin/gophermat/internal/model"
	"github.com/golangTroshin/gophermat/internal/repository"
)

var (
	ErrOrderExistsSameUser      = errors.New("order already uploaded by the same user")
	ErrOrderExistsDifferentUser = errors.New("order already uploaded by a different user")
	ErrInvalidOrderNumber       = errors.New("invalid order number")
)

type OrderService interface {
	ValidateOrderNumber(orderNumber string) bool
	CreateOrder(userID uint, orderNumber string) error
	GetOrdersByUserID(userID uint) ([]model.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) ValidateOrderNumber(orderNumber string) bool {
	return isValidLuhn(orderNumber)
}

func (s *orderService) CreateOrder(userID uint, orderNumber string) error {
	existingOrder, err := s.orderRepo.FindOrderByNumber(orderNumber)
	if err == nil {
		if existingOrder.UserID == userID {
			return ErrOrderExistsSameUser
		}
		return ErrOrderExistsDifferentUser
	}

	newOrder := &model.Order{
		Number:     orderNumber,
		UserID:     userID,
		UploadedAt: time.Now(),
	}
	return s.orderRepo.CreateOrder(newOrder)
}

func (s *orderService) GetOrdersByUserID(userID uint) ([]model.Order, error) {
	return s.orderRepo.FindOrdersByUserID(userID)
}

func isValidLuhn(orderNumber string) bool {
	var sum int
	var alt bool
	for i := len(orderNumber) - 1; i >= 0; i-- {
		n, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
