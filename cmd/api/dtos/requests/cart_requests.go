package requests

type AddCartItemRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Slug      string  `json:"slug" validate:"required"`
	Qty       int     `json:"qty" validate:"required,min=1"`
	Image     string  `json:"image"`
	Price     float64 `json:"price" validate:"required"`
}

type RemoveCartItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
}
