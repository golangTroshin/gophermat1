package model

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Login     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	BalanceID uint   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Balance   UserBalance `gorm:"foreignKey:BalanceID;references:ID"`
}
