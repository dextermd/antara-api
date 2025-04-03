package services

import (
	"antara-api/internal/models"
	"errors"
	"gorm.io/gorm"
)

type MCategoryService struct {
	db *gorm.DB
}

func NewMCategoryService(db *gorm.DB) *MCategoryService {
	return &MCategoryService{db: db}
}

func (c MCategoryService) List(db *gorm.DB) ([]models.MCategoryModel, error) {
	var categories []models.MCategoryModel
	result := db.Find(&categories)
	if result.Error != nil {
		return nil, errors.New("failed to fetch all categories")
	}

	return categories, nil
}

func (c MCategoryService) Create(db *gorm.DB) {

}
