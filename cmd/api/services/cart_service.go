package services

import (
	"antara-api/internal/models"
	"errors"
	"gorm.io/gorm"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

func (cartService *CartService) GetCart(sessionCartID string, userID *uint) (*models.CartModel, error) {
	var cart models.CartModel
	var result *gorm.DB

	baseQuery := cartService.db.Preload("Items")

	if userID != nil {
		sessionCart, err := cartService.updateSessionCartIfExists(sessionCartID, userID)
		if err != nil {
			return nil, err
		}

		if sessionCart != nil {
			return sessionCart, nil
		}

		result = baseQuery.Where("user_id = ?", *userID).First(&cart)
	} else {
		result = baseQuery.Where("session_cart_id = ?", sessionCartID).First(&cart)
	}

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &cart, nil
}

func (cartService *CartService) Create(cart *models.CartModel) error {
	result := cartService.db.Create(cart)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cartService *CartService) Update(cart *models.CartModel) error {
	result := cartService.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(cart)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cartService *CartService) updateSessionCartIfExists(sessionCartID string, userID *uint) (*models.CartModel, error) {
	var sessionCart models.CartModel

	result := cartService.db.Where("session_cart_id = ? AND user_id IS NULL", sessionCartID).First(&sessionCart)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	sessionCart.UserID = userID
	updateResult := cartService.db.Save(&sessionCart)
	if updateResult.Error != nil {
		return nil, updateResult.Error
	}

	var updatedCart models.CartModel
	finalResult := cartService.db.Preload("Items").Where("id = ?", sessionCart.ID).First(&updatedCart)
	if finalResult.Error != nil {
		return nil, finalResult.Error
	}

	return &updatedCart, nil
}
