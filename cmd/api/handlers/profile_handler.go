package handlers

import (
	"antara-api/cmd/api/requests"
	"antara-api/cmd/api/services"
	"antara-api/common"
	"antara-api/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAuthenticatedUser(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	return common.SendSuccessResponse(c, "Test Authenticated successfully", user)
}

func (h *Handler) ChangeUserPassword(c echo.Context) error {
	userService := services.NewUserService(h.DB)

	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	payload := new(requests.ChangePasswordRequest)
	if err := c.Bind(payload); err != nil {

		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	if common.ComparePasswordHash(user.Password, payload.CurrentPassword) == false {
		return common.SendBadRequestResponse(c, "the supplied password does not match the current password")
	}

	err := userService.ChangeUserPassword(&user, payload.Password)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Password changed successfully", nil)
}
