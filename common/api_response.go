package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ValidationError struct {
	Error     string `json:"error"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}

type JSONSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type JSONFailedValidationResponse struct {
	Success bool               `json:"success"`
	Errors  []*ValidationError `json:"errors"`
	Message string             `json:"message"`
}

type JSONErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SendSuccessResponse(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, JSONSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendFailedValidationResponse(c echo.Context, errors []*ValidationError) error {
	return c.JSON(http.StatusUnprocessableEntity, JSONFailedValidationResponse{
		Success: false,
		Errors:  errors,
		Message: "Validation failed",
	})
}

func SendErrorResponse(c echo.Context, message string, statusCode int) error {
	return c.JSON(statusCode, JSONErrorResponse{
		Success: false,
		Message: message,
	})
}

func SendBadRequestResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusBadRequest)
}

func SendNotFoundResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusNotFound)
}

func SendInternalServerErrorResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusInternalServerError)
}

func SendForbiddenResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusForbidden)
}

func SendUnauthorizedResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusUnauthorized)
}
