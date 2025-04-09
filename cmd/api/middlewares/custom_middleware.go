package middlewares

import (
	fnt "fmt"
	"github.com/labstack/echo/v4"
)

func CustomMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fnt.Println("we are in the custom middleware")
		return next(c)
	}
}
