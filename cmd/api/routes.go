package main

import (
	"antara-api/cmd/api/handlers"
	"github.com/labstack/echo/v4"
)

func (app *Application) routes(e *echo.Echo, handler handlers.Handler) {
	apiGroup := app.server.Group("/api")

	publicAuthRoutes := apiGroup.Group("/auth")
	{
		publicAuthRoutes.POST("/register", handler.SignUpHandler)
		publicAuthRoutes.POST("/login", handler.SignInHandler)
		publicAuthRoutes.POST("/forgot/password", handler.ForgotPasswordHandler)
		publicAuthRoutes.POST("/reset/password", handler.ResetPasswordHandler)
		publicAuthRoutes.POST("/refresh", handler.RefreshTokenHandler)
		publicAuthRoutes.POST("/logout", handler.LogoutHandler)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthenticationMiddleware)
	{
		profileRoutes.GET("/authenticated/user", handler.GetAuthenticatedUser)
		profileRoutes.PATCH("/change/password", handler.ChangeUserPassword)
	}

	mCategoryRoutes := apiGroup.Group("/m-categories", app.appMiddleware.AuthenticationMiddleware)
	{
		mCategoryRoutes.GET("/all", handler.ListCategoriesHandler)
	}

	productRoutes := apiGroup.Group("/products", app.appMiddleware.AuthenticationMiddleware)
	{
		productRoutes.GET("/all", handler.ListProductsHandler)
		productRoutes.GET("/:slug", handler.GetProductBySlagHandler)
	}

	apiGroup.GET("", handler.HealthCheck)
	apiGroup.GET("/test-authenticated", handler.TestAuthenticatedUser, app.appMiddleware.AuthenticationMiddleware)
}
