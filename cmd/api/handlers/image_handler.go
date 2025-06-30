package handlers

import (
	"antara-api/cmd/api/services"
	"antara-api/common"

	"github.com/labstack/echo/v4"
)

func (h *Handler) UploadTempImageHandler(c echo.Context) error {

	sessionID, ok := c.Get("session_id").(string)
	if !ok || sessionID == "" {
		return common.SendUnauthorizedResponse(c, "Session ID is required")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return common.SendBadRequestResponse(c, "No file provided")
	}

	if file.Size > 20*1024*1024 {
		return common.SendBadRequestResponse(c, "File size too large. Maximum 20MB allowed")
	}

	imageService := services.NewImageService(h.DB, "./uploads")
	uploadedImage, err := imageService.UploadTempImage(sessionID, file)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	response := map[string]string{
		"link": uploadedImage.ImageURL,
	}

	return c.JSON(200, response)
}

func (h *Handler) UploadPageImageHandler(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return common.SendBadRequestResponse(c, "Page slug is required")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return common.SendBadRequestResponse(c, "No file provided")
	}

	if file.Size > 10*1024*1024 {
		return common.SendBadRequestResponse(c, "File size too large. Maximum 10MB allowed")
	}

	imageService := services.NewImageService(h.DB, "./uploads")
	uploadedImage, err := imageService.UploadPageImage(slug, file)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	response := map[string]string{
		"link": uploadedImage.ImageURL,
	}

	return c.JSON(200, response)
}
