package handlers

import (
	"antara-api/cmd/api/requests"
	"antara-api/cmd/api/services"
	"antara-api/common"
	"antara-api/internal/mailer"
	"antara-api/internal/models"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) SignUpHandler(c echo.Context) error {
	payload := new(requests.SignUpRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	existingUser, err := userService.GetByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == false && existingUser != nil {
			return common.SendBadRequestResponse(c, "Email already exists")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	registeredUser, err := userService.CreateUser(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	mailData := mailer.EmailData{
		Subject: "Welcome to Antara",
		Meta: struct {
			Name      string
			LoginLink string
		}{
			Name:      *registeredUser.Name,
			LoginLink: "#",
		},
	}

	err = h.Mailer.SendMail(payload.Email, "welcome.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}

	return common.SendSuccessResponse(c, "User has been created", registeredUser)
}

func (h *Handler) SignInHandler(c echo.Context) error {
	userService := services.NewUserService(h.DB)

	payload := new(requests.SignInRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	userRetrieved, err := userService.GetByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && userRetrieved == nil {
			return common.SendBadRequestResponse(c, "Invalid email or password")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	if common.ComparePasswordHash(userRetrieved.Password, payload.Password) == false {
		return common.SendBadRequestResponse(c, "Invalid email or password")
	}

	accessToken, refreshToken, err := common.GenerateJWT(*userRetrieved)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "User logged in", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userRetrieved,
	})
}

func (h *Handler) TestAuthenticatedUser(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	return common.SendSuccessResponse(c, "Test Authenticated successfully", user)
}
