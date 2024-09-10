package model

import "time"

type UserBalance struct {
	ID        uint    `gorm:"primaryKey"`
	Current   float64 `gorm:"not null;default:0"`
	Withdrawn float64 `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
