package services

import (
	"antara-api/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageService struct {
	db         *gorm.DB
	uploadPath string
}

func NewImageService(db *gorm.DB, uploadPath string) *ImageService {
	return &ImageService{
		db:         db,
		uploadPath: uploadPath,
	}
}

func (s *ImageService) UploadTempImage(sessionID string, file *multipart.FileHeader) (*models.TempImageModel, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if !s.isValidImageType(file.Filename) {
		return nil, fmt.Errorf("invalid file type. Only jpg, jpeg, png, gif, webp are allowed")
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("temp_%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	tempDir := filepath.Join(s.uploadPath, "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(tempDir, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	imageURL := fmt.Sprintf("/uploads/temp/%s", newFileName)

	tempImage := &models.TempImageModel{
		ImageURL:  imageURL,
		FileName:  newFileName,
		FileSize:  file.Size,
		SessionID: sessionID,
	}

	if err := s.db.Create(tempImage).Error; err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save temp image info: %w", err)
	}

	return tempImage, nil
}

func (s *ImageService) UploadPageImage(slug string, file *multipart.FileHeader) (*models.PageImageModel, error) {
	var page models.PageModel
	if err := s.db.Where("slug = ?", slug).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("page with slug '%s' not found", slug)
		}
		return nil, fmt.Errorf("failed to find page: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if !s.isValidImageType(file.Filename) {
		return nil, fmt.Errorf("invalid file type. Only jpg, jpeg, png, gif, webp are allowed")
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	pageDir := filepath.Join(s.uploadPath, "pages", fmt.Sprintf("page_%s", slug))
	if err := os.MkdirAll(pageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(pageDir, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	imageURL := fmt.Sprintf("/uploads/pages/page_%s/%s", slug, newFileName)

	pageImage := &models.PageImageModel{
		PageID:   page.ID,
		ImageURL: imageURL,
		FileName: newFileName,
		FileSize: file.Size,
	}

	if err := s.db.Create(pageImage).Error; err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image info: %w", err)
	}

	return pageImage, nil
}

func (s *ImageService) MoveTempImagesToPage(slug string, PageID uint, content string, sessionID string) error {
	var tempImages []models.TempImageModel
	if err := s.db.Where("session_id = ?", sessionID).Find(&tempImages).Error; err != nil {
		return fmt.Errorf("failed to find temp images: %w", err)
	}

	if len(tempImages) == 0 {
		return nil
	}

	pageDir := filepath.Join(s.uploadPath, "pages", fmt.Sprintf("page_%s", slug))
	if err := os.MkdirAll(pageDir, 0755); err != nil {
		return fmt.Errorf("failed to create page directory: %w", err)
	}

	for _, tempImage := range tempImages {
		if !strings.Contains(content, tempImage.ImageURL) {
			s.deleteTempImage(tempImage)
			continue
		}

		ext := filepath.Ext(tempImage.FileName)
		newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

		oldPath := filepath.Join(s.uploadPath, "temp", tempImage.FileName)
		newPath := filepath.Join(pageDir, newFileName)

		if err := os.Rename(oldPath, newPath); err != nil {
			if err := s.copyFile(oldPath, newPath); err != nil {
				fmt.Printf("Warning: failed to move temp image %s: %v\n", tempImage.FileName, err)
				continue
			}
			os.Remove(oldPath)
		}

		newImageURL := fmt.Sprintf("/uploads/pages/page_%s/%s", slug, newFileName)
		pageImage := &models.PageImageModel{
			PageID:   PageID,
			ImageURL: newImageURL,
			FileName: newFileName,
			FileSize: tempImage.FileSize,
		}

		if err := s.db.Create(pageImage).Error; err != nil {
			fmt.Printf("Warning: failed to create page image record: %v\n", err)
			continue
		}

		content = strings.ReplaceAll(content, tempImage.ImageURL, newImageURL)

		s.db.Delete(&tempImage)
	}

	if err := s.db.Model(&models.PageModel{}).Where("id = ?", PageID).Update("content", content).Error; err != nil {
		return fmt.Errorf("failed to update page content: %w", err)
	}

	return nil
}

// Очистка старых временных изображений (запускать по cron)
func (s *ImageService) CleanupOldTempImages(olderThanHours int) error {
	cutoffTime := time.Now().Add(-time.Duration(olderThanHours) * time.Hour)

	var oldTempImages []models.TempImageModel
	if err := s.db.Where("created_at < ?", cutoffTime).Find(&oldTempImages).Error; err != nil {
		return fmt.Errorf("failed to find old temp images: %w", err)
	}

	for _, tempImage := range oldTempImages {
		s.deleteTempImage(tempImage)
	}

	if err := s.db.Where("created_at < ?", cutoffTime).Delete(&models.TempImageModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete old temp image records: %w", err)
	}

	return nil
}

func (s *ImageService) deleteTempImage(tempImage models.TempImageModel) {
	filePath := filepath.Join(s.uploadPath, "temp", tempImage.FileName)
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Warning: failed to delete temp file %s: %v\n", filePath, err)
	}
	s.db.Delete(&tempImage)
}

func (s *ImageService) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (s *ImageService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validType := range validTypes {
		if ext == validType {
			return true
		}
	}
	return false
}

func (s *ImageService) DeletePageImages(pageID uint) error {
	var page models.PageModel
	if err := s.db.Where("id = ?", pageID).First(&page).Error; err != nil {
		return fmt.Errorf("failed to find page: %w", err)
	}

	var pageImages []models.PageImageModel
	if err := s.db.Where("page_id = ?", pageID).Find(&pageImages).Error; err != nil {
		return fmt.Errorf("failed to find page images: %w", err)
	}

	for _, pageImage := range pageImages {
		relativePath := strings.TrimPrefix(pageImage.ImageURL, "/uploads/")
		filePath := filepath.Join(s.uploadPath, relativePath)

		fmt.Printf("Deleting file: %s\n", filePath)
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Warning: failed to delete page image file %s: %v\n", filePath, err)
		}
	}

	if err := s.db.Where("page_id = ?", pageID).Delete(&models.PageImageModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete page image records: %w", err)
	}

	pageDir := filepath.Join(s.uploadPath, "pages", fmt.Sprintf("page_%s", page.Slug))

	fmt.Printf("Attempting to delete directory: %s\n", pageDir)

	if _, err := os.Stat(pageDir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", pageDir)
		return nil
	}

	files, err := os.ReadDir(pageDir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", pageDir, err)
	} else {
		fmt.Printf("Directory %s contains %d items:\n", pageDir, len(files))
		for _, file := range files {
			fmt.Printf("  - %s (IsDir: %v)\n", file.Name(), file.IsDir())
		}
	}

	if err := os.RemoveAll(pageDir); err != nil {
		fmt.Printf("Warning: failed to remove page directory %s: %v\n", pageDir, err)
	}

	return nil
}
func (s *ImageService) DeletePageImagesBySlug(slug string) error {
	var page models.PageModel
	if err := s.db.Where("slug = ?", slug).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("failed to find page: %w", err)
	}

	return s.DeletePageImages(page.ID)
}
