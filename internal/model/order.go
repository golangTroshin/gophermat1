package model

import (
	"time"
)

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	Number     string    `gorm:"unique;not null"`
	UserID     uint      `gorm:"not null"`
	User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status     string    `gorm:"not null;default:NEW"`
	Accrual    float64   `gorm:"default:0.0"`
	UploadedAt time.Time `gorm:"not null;default:current_timestamp"`
}
