package repository

import (
	"github.com/golangTroshin/gophermat/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) error
	FindOrderByNumber(number string) (*model.Order, error)
	FindOrdersByUserID(userID uint) ([]model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindOrderByNumber(number string) (*model.Order, error) {
	var order model.Order
	if err := r.db.Where("number = ?", number).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindOrdersByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Where("user_id = ?", userID).Order("uploaded_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
