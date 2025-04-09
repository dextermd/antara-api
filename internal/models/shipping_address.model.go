package models

type ShippingAddressModel struct {
	BaseModel
	UserID     uint   `gorm:"not null" json:"user_id"`
	Phone      string `gorm:"type:varchar(255);"`
	Address    string `gorm:"type:varchar(255);"`
	City       string `gorm:"type:varchar(255);"`
	Country    string `gorm:"type:varchar(255);"`
	PostalCode string `gorm:"type:varchar(255);"`
}

func (u *ShippingAddressModel) TableName() string {
	return "shipping_addresses"
}
