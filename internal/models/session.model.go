package models

import (
	"time"
)

type SessionModel struct {
	SessionID    string    `gorm:"primaryKey;autoIncrement:false;type:varchar(255);index" json:"session_id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	AccessToken  string    `gorm:"type:text" json:"access_token"`
	RefreshToken string    `gorm:"type:text" json:"refresh_token"`
	IPAddress    string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent    string    `gorm:"type:text" json:"user_agent,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiresAt    time.Time `gorm:"index" json:"expires_at"`
}

func (SessionModel) TableName() string {
	return "sessions"
}

func (s SessionModel) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
