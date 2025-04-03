package main

import (
	"antara-api/cmd/api/handlers"
	"antara-api/cmd/api/middlewares"
	"antara-api/common"
	"antara-api/internal/mailer"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

type Application struct {
	logger        echo.Logger
	server        *echo.Echo
	handler       handlers.Handler
	appMiddleware middlewares.AppMiddleware
}

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	db, err := common.DBConnect()

	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	appMailer := mailer.NewMailer(e.Logger)
	h := handlers.Handler{
		DB:     db,
		Logger: e.Logger,
		Mailer: appMailer,
	}

	appMiddleware := middlewares.AppMiddleware{
		Logger: e.Logger,
		DB:     db,
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4200", "http://127.0.0.1:4200"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Cookie", "Set-Cookie"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
	}))

	app := Application{
		logger:        e.Logger,
		server:        e,
		handler:       h,
		appMiddleware: appMiddleware,
	}

	e.Use(middlewares.CustomMiddleware, middleware.Logger(), middleware.Recover())

	app.routes(e, h)
	fmt.Println(app)

	e.Static("/uploads", "./uploads")

	port := os.Getenv("APP_PORT")
	host := os.Getenv("APP_HOST")

	e.Logger.Fatal(e.Start(host + ":" + port))
}
