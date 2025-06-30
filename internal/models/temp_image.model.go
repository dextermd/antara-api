package models

import "time"

type TempImageModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ImageURL  string    `gorm:"type:varchar(500);not null" json:"image_url"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	SessionID string    `gorm:"type:varchar(255);index" json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TempImageModel) TableName() string {
	return "temp_images"
}
