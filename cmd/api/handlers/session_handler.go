package handlers

import (
	"antara-api/cmd/api/services"
	"antara-api/common"
	"antara-api/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetSessions(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	sessionID, ok := c.Get("session_id").(string)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "Session ID not found in context")
	}

	sessionService := services.NewSessionService(h.DB)
	sessions, err := sessionService.GetSessions(user.ID, sessionID)
	if err != nil {
		h.Logger.Error("Failed to retrieve sessions: ", err)
		return common.SendInternalServerErrorResponse(c, "Failed to retrieve sessions")
	}
	return common.SendSuccessResponse(c, "Ok", sessions)
}
