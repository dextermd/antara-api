package middlewares

import (
	"antara-api/common"
	"antara-api/internal/models"
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
		authHeader := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") == false {
			return common.SendUnauthorizedResponse(c, "Please provide a Bearer token")
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := common.ParseJWTSignedAccessToken(accessToken)
		if err != nil {
			return common.SendUnauthorizedResponse(c, "Invalid access token")
		}

		if common.IsClaimExpired(claims) {
			return common.SendUnauthorizedResponse(c, "Token expired")
		}

		var user models.UserModel
		result := appMiddleware.DB.First(&user, claims.ID)
		if result.Error != nil {
			return common.SendUnauthorizedResponse(c, "Invalid access token")
		}

		c.Set("user", user)

		return next(c)
	}
}
