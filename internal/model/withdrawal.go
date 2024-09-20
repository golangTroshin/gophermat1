package model

type UserWithdrawal struct {
	ID          uint    `gorm:"primaryKey"`
	UserID      uint    `gorm:"not null"`
	User        User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OrderNumber string  `gorm:"not null"`
	Sum         float64 `gorm:"not null"`
	ProcessedAt string  `gorm:"not null"`
}
