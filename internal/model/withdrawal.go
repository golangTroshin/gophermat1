package model

import (
	"time"
)

type UserWithdrawal struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	Order       string    `gorm:"not null"`
	Sum         float64   `gorm:"not null"`
	ProcessedAt time.Time `gorm:"not null"`
}
