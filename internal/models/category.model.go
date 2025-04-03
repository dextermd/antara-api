package models

type CategoryModel struct {
	BaseModel
	Name          string          `gorm:"type:varchar(60);not null" json:"name"`
	Slug          string          `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Description   *string         `gorm:"type:text" json:"description"`
	ParentId      *uint           `gorm:"index" json:"parent_id"`
	Parent        *CategoryModel  `gorm:"foreignKey:ParentId" json:"parent"`
	Subcategories []CategoryModel `gorm:"foreignKey:ParentId" json:"subcategories"`
	Products      []ProductModel  `gorm:"foreignKey:CategoryId" json:"products"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
