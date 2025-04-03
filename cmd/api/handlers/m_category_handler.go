package handlers

import (
	"antara-api/cmd/api/services"
	"antara-api/common"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ListCategoriesHandler(c echo.Context) error {
	mCategoryService := services.NewMCategoryService(h.DB)
	categories, err := mCategoryService.List(h.DB)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "ok", categories)
}
