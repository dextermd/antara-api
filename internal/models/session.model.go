package models

import (
	"gorm.io/gorm"
	"time"
)

type SessionModel struct {
	ID           string         `gorm:"primaryKey;type:varchar(255);not null" json:"id,omitempty"`
	UserID       uint           `gorm:"index" json:"user_id,omitempty"`
	User         UserModel      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device       string         `gorm:"type:varchar(255)" json:"device,omitempty"`
	UserAgent    string         `gorm:"type:text" json:"user_agent,omitempty"`
	IPAddress    string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	LastActivity time.Time      `gorm:"index" json:"last_activity,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active,omitempty"`
	ExpiresAt    time.Time      `gorm:"index" json:"expires_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SessionModel) TableName() string {
	return "sessions"
}
