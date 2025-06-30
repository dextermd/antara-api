package models

import "time"

type PageImageModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PageID    uint      `gorm:"index;not null" json:"page_id"`
	ImageURL  string    `gorm:"type:varchar(500);not null" json:"image_url"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Page PageModel `gorm:"foreignKey:PageID;constraint:OnDelete:CASCADE" json:"page,omitempty"`
}

func (PageImageModel) TableName() string {
	return "page_images"
}
