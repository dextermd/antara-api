package handlers

import (
	"antara-api/cmd/api/dtos/requests"
	"antara-api/cmd/api/dtos/response"
	"antara-api/cmd/api/services"
	"antara-api/common"
	"antara-api/internal/mailer"
	"antara-api/internal/models"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
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
			Name:      *registeredUser.FirstName,
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
	payload := new(requests.SignInRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	userRetrieved, err := userService.GetByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && userRetrieved == nil {
			return common.SendBadRequestResponse(c, "Invalid email or password")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	if userRetrieved.IsActive == false {
		return common.SendBadRequestResponse(c, "User is not active")
	}

	if common.ComparePasswordHash(userRetrieved.PasswordHash, payload.Password) == false {
		return common.SendBadRequestResponse(c, "Invalid email or password")
	}

	sessionService := services.NewSessionService(h.DB)
	session, err := sessionService.CreateSession(userRetrieved.ID, payload.Device, c.Request().Header.Get("User-Agent"), c.RealIP())
	if err != nil {
		h.Logger.Error("Failed to create session: ", err)
		return common.SendInternalServerErrorResponse(c, "Failed to create session")
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true для HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	}
	c.SetCookie(cookie)

	return common.SendSuccessResponse(c, "User logged in",
		&response.AuthDataResponse{
			User:      userRetrieved,
			SessionID: session.ID,
			ExpiresAt: session.ExpiresAt,
		})
}

func (h *Handler) TestAuthenticatedUser(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	return common.SendSuccessResponse(c, "Test Authenticated successfully", user)
}

func (h *Handler) LogoutHandler(c echo.Context) error {
	sessionID, _ := c.Get("session_id").(string)

	fmt.Println("Session ID from context:", sessionID)

	sessionService := services.NewSessionService(h.DB)
	err := sessionService.RevokeSession(sessionID)
	if err != nil {
		fmt.Println("Failed to revoke session:", err)
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)

	return common.SendSuccessResponse(c, "User logged out", nil)
}
