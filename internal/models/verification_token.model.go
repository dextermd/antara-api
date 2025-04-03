package models

import "time"

type VerificationTokenModel struct {
	Token      string    `json:"-" gorm:"primaryKey;autoIncrement:false;type:varchar(255);not null"`
	TargetId   uint      `json:"-" gorm:"index;not null"`
	Identifier string    `json:"-" gorm:"primaryKey;autoIncrement:false;type:varchar(255);not null"`
	Used       bool      `json:"-" gorm:"index;type:bool"`
	ExpiresAt  time.Time `json:"-" gorm:"index;not null"`
}

func (VerificationTokenModel) TableName() string {
	return "verification_tokens"
}
