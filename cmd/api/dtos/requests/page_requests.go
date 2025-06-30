package requests

type CreatePageRequest struct {
	Title           string `json:"title" validate:"required,min=1,max=255"`
	Slug            string `json:"slug" validate:"required,min=1,max=255"`
	Content         string `json:"content" validate:"required"`
	IsPublished     bool   `json:"is_published" validate:"required"`
	MetaTitle       string `json:"meta_title" validate:"max=255"`
	MetaDescription string `json:"meta_description" validate:"max=255"`
	MetaKeywords    string `json:"meta_keywords" validate:"max=255"`
	DisplayOrder    int    `json:"display_order" validate:"gte=0"`
	PageType        string `json:"page_type" validate:"required,min=1,max=255"`
	RoutePath       string `json:"route_path" validate:"required,min=1,max=255"`
}

type UpdatePageRequest struct {
	Title           string `json:"title" validate:"omitempty,min=1,max=255"`
	Slug            string `json:"slug" validate:"omitempty,min=1,max=255"`
	Content         string `json:"content" validate:"omitempty"`
	IsPublished     *bool  `json:"is_published" validate:"omitempty"`
	MetaTitle       string `json:"meta_title" validate:"omitempty,max=255"`
	MetaDescription string `json:"meta_description" validate:"omitempty,max=255"`
	MetaKeywords    string `json:"meta_keywords" validate:"omitempty,max=255"`
	DisplayOrder    *int   `json:"display_order" validate:"omitempty,gte=0"`
	PageType        string `json:"page_type" validate:"omitempty,min=1,max=255"`
	RoutePath       string `json:"route_path" validate:"omitempty,min=1,max=255"`
}
