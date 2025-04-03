package models

type OptionModel struct {
	BaseModel
	Name      string       `gorm:"type:varchar(50);not null" json:"name"`
	Values    []string     `gorm:"type:text[]" json:"values"`
	Price     float64      `gorm:"type:decimal(12,2);default:0.0" json:"price"`
	ProductId uint         `gorm:"index" json:"productId"`
	Product   ProductModel `gorm:"foreignKey:ProductId" json:"product"`
}

func (OptionModel) TableName() string {
	return "options"
}
