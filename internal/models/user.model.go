package models

import (
	"time"
)

type UserModel struct {
	BaseModel
	FirstName     *string               `gorm:"type:varchar(255)" json:"first_name,omitempty"`
	LastName      *string               `gorm:"type:varchar(255)" json:"last_name,omitempty"`
	Email         string                `gorm:"type:varchar(255);unique;not null" json:"email,omitempty"`
	EmailVerified *time.Time            `gorm:"type:timestamp;default:null" json:"email_verified,omitempty"`
	AvatarURL     *string               `gorm:"type:text" json:"avatar_url,omitempty"`
	PasswordHash  string                `gorm:"type:varchar(255);default:null" json:"-"`
	IsActive      bool                  `gorm:"default:true" json:"is_active,omitempty"`
	Sessions      []SessionModel        `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	Providers     []SocialProviderModel `gorm:"foreignKey:UserID" json:"providers,omitempty"`
	Roles         []RoleModel           `gorm:"many2many:user_has_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles,omitempty"`
}

func (u *UserModel) TableName() string {
	return "users"
}
