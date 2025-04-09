package models

type CartModel struct {
	BaseModel
	UserID        *uint           `gorm:"index;null" json:"user_id"`
	SessionCartID string          `gorm:"type:varchar(255);not null" json:"session_cart_id"`
	Items         []CartItemModel `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	ItemsPrice    float64         `gorm:"type:decimal(12,2);not null"`
	TotalJustItem float64         `gorm:"type:decimal(12,2);not null"`
	TotalPrice    float64         `gorm:"type:decimal(12,2);not null"`
	ShippingPrice float64         `gorm:"type:decimal(12,2);not null"`
	TaxPrice      float64         `gorm:"type:decimal(12,2);not null"`
}

func (u *CartModel) TableName() string {
	return "carts"
}
