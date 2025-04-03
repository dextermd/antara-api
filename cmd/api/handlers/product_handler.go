package handlers

import (
	"antara-api/cmd/api/services"
	"antara-api/common"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *Handler) ListProductsHandler(c echo.Context) error {
	productService := services.NewProductService(h.DB)
	products, err := productService.List(h.DB)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "Authentication",
		Value:    "*accessToken",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	c.SetCookie(&http.Cookie{
		Name:     "Refresh",
		Value:    "*refreshToken",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(100 * time.Hour),
	})

	return common.SendSuccessResponse(c, "ok", products)
}

func (h *Handler) GetProductBySlagHandler(c echo.Context) error {
	fmt.Println(c.Param("slug"))
	productService := services.NewProductService(h.DB)
	product, err := productService.GetBySlug(h.DB, c.Param("slug"))
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "ok", product)
}
