package db

import (
	"github.com/golangTroshin/gophermat/internal/config"
	"github.com/golangTroshin/gophermat/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.Options.DataBaseURI), &gorm.Config{})
	if err != nil {
		return db, err
	}

	err = db.AutoMigrate(&model.User{}, &model.UserBalance{}, &model.Order{}, &model.UserWithdrawal{})
	if err != nil {
		return db, err
	}

	return db, nil
}
