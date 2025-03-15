package models

import "time"

type AppTokenModel struct {
	BaseModel
	Token     string    `json:"-" gorm:"index;type:varchar(255);not null"`
	TargetId  uint      `json:"-" gorm:"index;not null"`
	Type      string    `json:"-" gorm:"index;type:varchar(255);not null"`
	Used      bool      `json:"-" gorm:"index;type:bool"`
	ExpiresAt time.Time `json:"-" gorm:"index;not null"`
}

func (AppTokenModel) TableName() string {
	return "app_tokens"
}
