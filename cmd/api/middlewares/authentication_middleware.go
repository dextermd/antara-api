package middlewares

import (
	"antara-api/cmd/api/services"
	"antara-api/common"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		cookie, err := c.Cookie("session_id")
		if err != nil {
			return common.SendUnauthorizedResponse(c, "Session ID not found in cookie or header")
		}

		sessionID := cookie.Value

		sessionService := services.NewSessionService(appMiddleware.DB)
		user, err := sessionService.ValidateSession(sessionID)
		if err != nil {
			fmt.Println("Session validation error:", err)
			return common.SendUnauthorizedResponse(c, err.Error())
		}

		c.Set("user", user)
		c.Set("session_id", sessionID)
		return next(c)
	}
}
