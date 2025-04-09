package services

import (
	"antara-api/internal/models"
	"errors"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (c ProductService) List(db *gorm.DB) ([]models.ProductModel, error) {
	var products []models.ProductModel
	result := db.Find(&products)
	if result.Error != nil {
		return nil, errors.New("failed to fetch all categories")
	}

	return products, nil
}

func (c ProductService) GetBySlug(db *gorm.DB, slug string) (*models.ProductModel, error) {
	var product models.ProductModel
	result := db.Where("slug = ?", slug).First(&product)
	if result.Error != nil {
		return nil, errors.New("failed to fetch product by slug")
	}

	return &product, nil
}

func (c ProductService) GetByID(db *gorm.DB, productID uint) (*models.ProductModel, error) {
	var product models.ProductModel
	result := db.Where("id = ?", productID).First(&product)
	if result.Error != nil {
		return nil, errors.New("failed to fetch product by slug")
	}

	return &product, nil
}
