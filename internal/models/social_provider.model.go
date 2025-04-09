package models

import (
	"time"
)

type SocialProviderModel struct {
	ID         uint      `gorm:"primaryKey" json:"id,omitempty"`
	UserID     *uint     `gorm:"index;null" json:"user_id,omitempty"`
	Provider   string    `gorm:"type:varchar(255);unique;not null" json:"provider,omitempty"`
	ProviderID string    `gorm:"type:varchar(255);not null" json:"provider_id,omitempty"`
	Email      string    `gorm:"type:varchar(255);not null" json:"email,omitempty"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	AvatarURL  string    `gorm:"type:text" json:"avatar_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (SocialProviderModel) TableName() string {
	return "social_providers"
}
