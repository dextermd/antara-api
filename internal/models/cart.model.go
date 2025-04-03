package models

import (
	"gorm.io/datatypes"
)

type CartModel struct {
	BaseModel
	UserID        *uint          `gorm:"index;null" json:"user_id"`
	SessionCartID string         `gorm:"type:varchar(255);not null"`
	Items         datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	ItemsPrice    float64        `gorm:"type:decimal(12,2);not null"`
	TotalJustItem float64        `gorm:"type:decimal(12,2);not null"`
	TotalPrice    float64        `gorm:"type:decimal(12,2);not null"`
	ShippingPrice float64        `gorm:"type:decimal(12,2);not null"`
	TaxPrice      float64        `gorm:"type:decimal(12,2);not null"`
	User          *UserModel     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (u *CartModel) TableName() string {
	return "carts"
}
