package models

import (
	"time"
)

type UserModel struct {
	BaseModel
	Name          *string    `gorm:"type:varchar(255);default:'No Name'" json:"name"`
	Email         string     `gorm:"type:varchar(255);unique;not null" json:"email"`
	EmailVerified *time.Time `gorm:"type:timestamp;default:null" json:"email_verified"`
	Images        *string    `gorm:"type:text;default:null" json:"images"`
	Password      string     `gorm:"type:varchar(255);default:null" json:"-"`
	PaymentMethod *string    `gorm:"type:varchar(255);default:null" json:"payment_method"`
	LastOnline    time.Time  `gorm:"autoUpdateTime" json:"last_online"`

	Roles             []RoleModel            `gorm:"many2many:user_has_roles;" json:"roles"`
	Carts             []CartModel            `gorm:"foreignKey:UserID" json:"carts"`
	ShippingAddresses []ShippingAddressModel `gorm:"foreignKey:UserID" json:"shipping_addresses"`
}

func (u *UserModel) TableName() string {
	return "users"
}
