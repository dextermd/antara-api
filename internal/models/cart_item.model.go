package models

type CartItemModel struct {
	BaseModel
	CartID    uint         `gorm:"index;not null" json:"cart_id"`
	ProductID uint         `gorm:"not null" json:"product_id"`
	Name      string       `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string       `gorm:"type:varchar(255);not null;unique" json:"slug"`
	Qty       int          `gorm:"not null" json:"qty"`
	Image     string       `gorm:"type:varchar(512)" json:"image"`
	Price     float64      `gorm:"type:decimal(12,2);not null" json:"price"`
	Cart      CartModel    `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Product   ProductModel `gorm:"foreignKey:ProductID"`
}

func (u *CartItemModel) TableName() string {
	return "cart_items"
}
