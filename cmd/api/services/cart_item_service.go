package services

import (
	"antara-api/internal/models"
	"gorm.io/gorm"
)

type CartItemService struct {
	db *gorm.DB
}

func NewCartItemService(db *gorm.DB) *CartItemService {
	return &CartItemService{db: db}
}

func (cartItemService *CartItemService) Update(cartItem *models.CartItemModel) error {
	result := cartItemService.db.Save(cartItem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cartItemService *CartItemService) Delete(cartItem *models.CartItemModel) error {
	result := cartItemService.db.Unscoped().Delete(cartItem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
