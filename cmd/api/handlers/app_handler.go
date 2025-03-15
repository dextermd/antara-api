package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type healthCheckStruct struct {
	Health bool `json:"health"`
}

func (h *Handler) HealthCheck(c echo.Context) error {
	healthCheckStruct := healthCheckStruct{
		Health: true,
	}
	return c.JSON(http.StatusOK, healthCheckStruct)
}
