package models

import "github.com/lib/pq"

type ProductModel struct {
	BaseModel
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Description *string        `gorm:"type:text" json:"description"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images" `
	Brand       string         `gorm:"type:varchar(255);not null" json:"brand"`
	Stock       int            `gorm:"type:int;default:0" json:"stock"`
	Price       float64        `gorm:"type:decimal(12,2);default:0.0" json:"price"`
	Rating      float64        `gorm:"type:decimal(4,2);default:0" json:"rating"`
	NumReviews  int            `gorm:"type:int;default:0" json:"num_reviews"`
	IsFeatured  bool           `gorm:"type:boolean;default:false" json:"is_featured"`
	Banner      *string        `gorm:"type:text" json:"banner"`
	CategoryId  *uint          `gorm:"index" json:"categoryId"`
	Category    *CategoryModel `gorm:"foreignKey:CategoryId" json:"category"`
	Options     []OptionModel  `gorm:"foreignKey:ProductId" json:"options"`
}

func (ProductModel) TableName() string {
	return "products"
}
