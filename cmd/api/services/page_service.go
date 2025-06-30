package services

import (
	"antara-api/cmd/api/dtos"
	"antara-api/cmd/api/dtos/requests"
	"antara-api/internal/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type PageService struct {
	db *gorm.DB
}

func NewPageService(db *gorm.DB) *PageService {
	return &PageService{db: db}
}

func (s *PageService) CreatePage(createPageRequest *requests.CreatePageRequest, sessionID string) (*models.PageModel, error) {
	page := &models.PageModel{
		Title:           createPageRequest.Title,
		Slug:            createPageRequest.Slug,
		Content:         createPageRequest.Content,
		IsPublished:     createPageRequest.IsPublished,
		MetaTitle:       createPageRequest.MetaTitle,
		MetaDescription: createPageRequest.MetaDescription,
		MetaKeywords:    createPageRequest.MetaKeywords,
		DisplayOrder:    createPageRequest.DisplayOrder,
		PageType:        createPageRequest.PageType,
		RoutePath:       createPageRequest.RoutePath,
	}

	result := s.db.Create(page)
	if result.Error != nil {
		return nil, result.Error
	}

	imageService := NewImageService(s.db, "./uploads")
	if err := imageService.MoveTempImagesToPage(page.Slug, page.ID, page.Content, sessionID); err != nil {
		fmt.Printf("Warning: failed to move temp images: %v\n", err)
	}

	return page, nil
}

func (s *PageService) ListPages(params *dtos.PaginationParams) ([]models.PageModel, int64, error) {
	var pages []models.PageModel
	var total int64

	query := s.db.Model(&models.PageModel{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.SortBy != "" {
		orderClause := params.SortBy
		if params.Order != "" {
			orderClause += " " + strings.ToUpper(params.Order)
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order("created_at DESC")
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&pages).Error; err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

func (s *PageService) GetPageBySlug(slug string) (*models.PageModel, error) {
	page := &models.PageModel{}
	result := s.db.Where("slug = ?", slug).First(page)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return page, nil
}

func (s *PageService) UpdatePage(slug string, updatePageRequest *requests.UpdatePageRequest) (*models.PageModel, error) {
	page := &models.PageModel{}
	result := s.db.Where("slug = ?", slug).First(page)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	if updatePageRequest.Title != "" {
		page.Title = updatePageRequest.Title
	}
	if updatePageRequest.Slug != "" {
		page.Slug = updatePageRequest.Slug
	}
	if updatePageRequest.Content != "" {
		page.Content = updatePageRequest.Content
	}
	if updatePageRequest.IsPublished != nil {
		page.IsPublished = *updatePageRequest.IsPublished
	}
	if updatePageRequest.MetaTitle != "" {
		page.MetaTitle = updatePageRequest.MetaTitle
	}
	if updatePageRequest.MetaDescription != "" {
		page.MetaDescription = updatePageRequest.MetaDescription
	}
	if updatePageRequest.MetaKeywords != "" {
		page.MetaKeywords = updatePageRequest.MetaKeywords
	}
	if updatePageRequest.DisplayOrder != nil {
		page.DisplayOrder = *updatePageRequest.DisplayOrder
	}
	if updatePageRequest.PageType != "" {
		page.PageType = updatePageRequest.PageType
	}
	if updatePageRequest.RoutePath != "" {
		page.RoutePath = updatePageRequest.RoutePath
	}

	result = s.db.Save(page)
	if result.Error != nil {
		return nil, result.Error
	}

	return page, nil
}

func (s *PageService) DeletePage(slug string) error {
	page := &models.PageModel{}
	result := s.db.Where("slug = ?", slug).First(page)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return result.Error
	}

	imageService := NewImageService(s.db, "./uploads")
	if err := imageService.DeletePageImages(page.ID); err != nil {
		fmt.Printf("Warning: failed to delete page images: %v\n", err)
	}

	result = s.db.Delete(page)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *PageService) ListPublishedPages() ([]models.PageModel, error) {
	var pages []models.PageModel
	result := s.db.Where("is_published = ?", true).Find(&pages)

	if result.Error != nil {
		return nil, result.Error
	}

	return pages, nil
}

func (s *PageService) GetPageImages(slug string) ([]models.PageImageModel, error) {
	page := &models.PageModel{}
	result := s.db.Where("slug = ?", slug).First(page)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	var images []models.PageImageModel
	result = s.db.Where("page_id = ?", page.ID).Find(&images)

	if result.Error != nil {
		return nil, result.Error
	}

	return images, nil
}
