package handlers

import (
	"antara-api/cmd/api/requests"
	"antara-api/cmd/api/services"
	"antara-api/common"
	"antara-api/internal/mailer"
	"encoding/base64"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/url"
)

func (h *Handler) ForgotPasswordHandler(c echo.Context) error {
	payload := new(requests.ForgotPasswordRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	existingUser, err := userService.GetByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == true && existingUser == nil {
			return common.SendNotFoundResponse(c, "Email does not exist, register with this email")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	token, err := appTokenService.GenerateResetPasswordToken(*existingUser)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(existingUser.Email))

	frontendUrl, err := url.Parse(payload.FrontendURL)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid frontend URL")
	}

	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", token.Token)
	frontendUrl.RawQuery = query.Encode()

	mailData := mailer.EmailData{
		Subject: "Request Password Reset",
		Meta: struct {
			Token       string
			FrontendUrl string
		}{
			Token:       token.Token,
			FrontendUrl: frontendUrl.String(),
		},
	}

	err = h.Mailer.SendMail(payload.Email, "password-reset.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}

	return common.SendSuccessResponse(c, "Password reset link sent to your email", nil)
}

func (h *Handler) ResetPasswordHandler(c echo.Context) error {
	payload := new(requests.ResetPasswordRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	email, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		return common.SendBadRequestResponse(c, "An error occurred, try again later")
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	existingUser, err := userService.GetByEmail(string(email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == true && existingUser == nil {
			return common.SendNotFoundResponse(c, "Invalid password reset token")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	appToken, err := appTokenService.ValidateResetPasswordToken(*existingUser, payload.Token)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	err = userService.ChangeUserPassword(existingUser, payload.Password)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	appTokenService.InvalidateToken(existingUser.Id, *appToken)

	return common.SendSuccessResponse(c, "Password reset successfully", nil)
}
