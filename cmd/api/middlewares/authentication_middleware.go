package middlewares

import (
	"antara-api/common"
	"antara-api/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strings"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")

		sessionCookie, err := c.Cookie("session_id")
		if err != nil {
			return common.SendUnauthorizedResponse(c, "No session found")
		}
		sessionID := sessionCookie.Value
		fmt.Println("Auth Middleware Session ID:", sessionID)

		var session models.SessionModel
		result := appMiddleware.DB.First(&session, "session_id = ?", sessionID)
		if result.Error != nil {
			return common.SendUnauthorizedResponse(c, "Session not found")
		}

		claims, err := common.ParseJWTSignedAccessToken(session.AccessToken)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				return common.SendUnauthorizedResponse(c, "Token is expired")
			}
			return common.SendUnauthorizedResponse(c, "Invalid access token")
		}

		var user models.UserModel
		data := appMiddleware.DB.First(&user, claims.ID)
		if data.Error != nil {
			return common.SendUnauthorizedResponse(c, "Invalid access token")
		}

		//// get cookie from request
		//cookie, err := c.Cookie("Authentication")
		//if err != nil {
		//	return common.SendUnauthorizedResponse(c, "Invalid access token")
		//}
		//accessToken := cookie.Value
		//fmt.Println("accessToken: ", accessToken)

		//authHeader := c.Request().Header.Get("Authorization")
		//if strings.HasPrefix(authHeader, "Bearer ") == false {
		//	return common.SendUnauthorizedResponse(c, "Please provide a Bearer token")
		//}
		//accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		//claims, err := common.ParseJWTSignedAccessToken(accessToken)
		//if err != nil {
		//	if strings.Contains(err.Error(), "token is expired") {
		//		return common.SendUnauthorizedResponse(c, "Token is expired")
		//	}
		//	return common.SendUnauthorizedResponse(c, "Invalid access token")
		//}
		//
		//var user models.UserModel
		//result := appMiddleware.DB.First(&user, claims.ID)
		//if result.Error != nil {
		//	return common.SendUnauthorizedResponse(c, "Invalid access token")
		//}

		c.Set("user", user)

		return next(c)
	}
}
