package models

type PageModel struct {
	BaseModel
	Title           string `gorm:"type:varchar(255);not null" json:"title"`
	Slug            string `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Content         string `gorm:"type:text" json:"content"`
	IsPublished     bool   `gorm:"default:false" json:"is_published"`
	MetaTitle       string `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDescription string `gorm:"type:varchar(255)" json:"meta_description"`
	MetaKeywords    string `gorm:"type:varchar(255)" json:"meta_keywords"`
	DisplayOrder    int    `gorm:"default:0" json:"display_order"`
	PageType        string `gorm:"type:varchar(255)" json:"page_type"`
	RoutePath       string `gorm:"type:varchar(255)" json:"route_path"`
}

func (PageModel) TableName() string {
	return "pages"
}
