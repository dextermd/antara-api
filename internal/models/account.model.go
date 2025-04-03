package models

import (
	"gorm.io/gorm"
	"time"
)

type AccountModel struct {
	UserID            uint           `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	Provider          string         `gorm:"type:varchar(255)" json:"provider"`
	ProviderAccountID string         `gorm:"primaryKey;autoIncrement:false;type:varchar(255)" json:"provider_account_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	User              UserModel      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

func (AccountModel) TableName() string {
	return "accounts"
}
