package handlers

import (
	"antara-api/cmd/api/dtos"
	"antara-api/cmd/api/dtos/requests"
	"antara-api/cmd/api/services"
	"antara-api/common"
	"math"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreatePageHandler(c echo.Context) error {
	payload := new(requests.CreatePageRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	sessionID, ok := c.Get("session_id").(string)
	if !ok || sessionID == "" {
		return common.SendUnauthorizedResponse(c, "Session ID is required")
	}

	pageService := services.NewPageService(h.DB)
	createdPage, err := pageService.CreatePage(payload, sessionID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Page has been created successfully", createdPage)
}

func (h *Handler) ListPagesHandler(c echo.Context) error {
	params := &dtos.PaginationParams{
		Page:     1,
		PageSize: 10,
		Order:    "asc",
	}

	if err := c.Bind(params); err != nil {
		return common.SendBadRequestResponse(c, "Invalid parameters")
	}

	pageService := services.NewPageService(h.DB)
	pages, total, err := pageService.ListPages(params)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	return common.SendSuccessPaginationResponse(c, "List of pages", pages, total, params.Page, params.PageSize, totalPages)
}

func (h *Handler) GetPageBySlugHandler(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return common.SendBadRequestResponse(c, "Slug is required")
	}

	pageService := services.NewPageService(h.DB)
	page, err := pageService.GetPageBySlug(slug)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	if page == nil {
		return common.SendNotFoundResponse(c, "Page not found")
	}

	return common.SendSuccessResponse(c, "Page retrieved successfully", page)
}

func (h *Handler) UpdatePageHandler(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return common.SendBadRequestResponse(c, "Slug is required")
	}

	payload := new(requests.UpdatePageRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	pageService := services.NewPageService(h.DB)
	updatedPage, err := pageService.UpdatePage(slug, payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Page has been updated successfully", updatedPage)
}

func (h *Handler) DeletePageHandler(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return common.SendBadRequestResponse(c, "Slug is required")
	}

	pageService := services.NewPageService(h.DB)
	err := pageService.DeletePage(slug)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Page has been deleted successfully", nil)
}

func (h *Handler) ListPublishedPagesHandler(c echo.Context) error {
	pageService := services.NewPageService(h.DB)
	pages, err := pageService.ListPublishedPages()
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "List of published pages", pages)
}

func (h *Handler) GetPageImagesHandler(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return common.SendBadRequestResponse(c, "Slug is required")
	}

	pageService := services.NewPageService(h.DB)
	images, err := pageService.GetPageImages(slug)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	if images == nil {
		return common.SendNotFoundResponse(c, "No images found for this page")
	}

	var result []map[string]string
	for _, img := range images {
		result = append(result, map[string]string{
			"url":   `http://localhost:8000` + img.ImageURL,
			"thumb": `http://localhost:8000` + img.ImageURL,
			"tag":   img.FileName,
		})
	}

	return c.JSON(200, result)
}
