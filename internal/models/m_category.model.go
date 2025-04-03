package models

type MCategoryModel struct {
	Name     string `gorm:"unique;type:varchar(200);not null" json:"name"`
	Slug     string `gorm:"unique;type:varchar(200);not null" json:"slug"`
	IsCustom bool   `gorm:"type:boolean;default:false" json:"is_custom"`
}

func (MCategoryModel) TableName() string {
	return "m_categories"
}
