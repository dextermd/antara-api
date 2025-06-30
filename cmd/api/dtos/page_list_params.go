package dtos

type PaginationParams struct {
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"pageSize" validate:"min=1,max=100"`
	Search   string `query:"search"`
	SortBy   string `query:"sortBy"`
	Order    string `query:"order" validate:"oneof=asc desc"`
}
