package handlers

import (
	"antara-api/cmd/api/services"
	"antara-api/common"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ListProductsHandler(c echo.Context) error {
	productService := services.NewProductService(h.DB)
	products, err := productService.List(h.DB)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "ok", products)
}

func (h *Handler) GetProductBySlagHandler(c echo.Context) error {
	productService := services.NewProductService(h.DB)
	product, err := productService.GetBySlug(h.DB, c.Param("slug"))
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "ok", product)
}
