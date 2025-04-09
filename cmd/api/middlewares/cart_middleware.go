package middlewares

import (
	"antara-api/common"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func CartMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var cartID string
		existingCookie, err := c.Cookie("session_cart_id")
		if err != nil || existingCookie.Value == "" {
			cartID = common.GenerateCartID()
			cookie := &http.Cookie{
				Name:     "session_cart_id",
				Value:    cartID,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().AddDate(0, 0, 30),
			}
			c.SetCookie(cookie)
		}

		return next(c)
	}
}
