package handlers

import (
	"antara-api/cmd/api/requests"
	"antara-api/cmd/api/services"
	"antara-api/common"
	"antara-api/internal/mailer"
	"antara-api/internal/models"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
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

	sessionID, err := common.GenerateSessionID()
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	session := models.SessionModel{
		UserID:       userRetrieved.ID,
		SessionID:    sessionID,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		IPAddress:    c.RealIP(),
		UserAgent:    c.Request().Header.Get("User-Agent"),
		ExpiresAt:    common.GetRefreshTokenExpirationTime(),
	}

	sessionService := services.NewSessionService(h.DB)
	_, err = sessionService.CreateSession(session)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return common.SendSuccessResponse(c, "User logged in", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userRetrieved,
	})
}

func (h *Handler) RefreshTokenHandler(c echo.Context) error {
	sessionCookie, err := c.Cookie("session_id")
	if err != nil {
		return common.SendBadRequestResponse(c, "No session found")
	}
	fmt.Println("RefreshTokenHandler -> Session ID: ", sessionCookie.Value)

	sessionService := services.NewSessionService(h.DB)
	session, err := sessionService.GetByID(sessionCookie.Value)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "Session not found")
	}

	claims, err := common.ParseJWTSignedRefreshToken(session.RefreshToken)
	sessionService.InvalidateSession(c, session.SessionID)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return common.SendForbiddenResponse(c, "Token is expired")
		}
		return common.SendForbiddenResponse(c, "Invalid access token")
	}

	userService := services.NewUserService(h.DB)
	user, err := userService.GetById(claims.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "User not found")
	}

	accessToken, refreshToken, err := common.GenerateJWT(*user)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	newSessionID, err := common.GenerateSessionID()
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	newSession := models.SessionModel{
		UserID:       user.ID,
		SessionID:    newSessionID,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		IPAddress:    c.RealIP(),
		UserAgent:    c.Request().Header.Get("User-Agent"),
		ExpiresAt:    common.GetRefreshTokenExpirationTime(),
	}

	_, err = sessionService.CreateSession(newSession)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return common.SendSuccessResponse(c, "Token refreshed", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
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
	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
	})
	sessionCookie, err := c.Cookie("session_id")
	if err == nil {
		sessionService := services.NewSessionService(h.DB)
		session, err := sessionService.GetByID(sessionCookie.Value)
		if err == nil {
			_ = sessionService.DeleteSession(session.SessionID)
		}
	}

	fmt.Println("LogoutHandler", sessionCookie)
	return common.SendSuccessResponse(c, "User logged out", nil)
}
