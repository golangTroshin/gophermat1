package repository

import (
	"github.com/golangTroshin/gophermat/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) CreateOrder(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindOrderByNumber(number string) (*model.Order, error) {
	var order model.Order
	if err := r.db.Where("number = ?", number).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindOrdersByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Where("user_id = ?", userID).Order("uploaded_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
