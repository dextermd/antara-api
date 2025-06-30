package cron

import (
	"antara-api/cmd/api/services"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func StartCleanupJob(db *gorm.DB) {
	imageService := services.NewImageService(db, "./uploads")

	ticker := time.NewTicker(24 * time.Hour) // Запускаем каждые 24 часа

	go func() {
		for range ticker.C {
			fmt.Println("Starting temp images cleanup job...")
			// Удаляем файлы старше 24 часов
			if err := imageService.CleanupOldTempImages(24); err != nil {
				fmt.Printf("Error during temp images cleanup: %v", err)
			} else {
				fmt.Println("Temp images cleanup completed successfully")
			}
		}
	}()
}
