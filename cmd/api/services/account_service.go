package services

import (
	"antara-api/internal/models"
	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (c AccountService) List(db *gorm.DB) ([]models.AccountModel, error) {
	return nil, nil
}
