package services

import (
	"antara-api/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (sessionService *SessionService) List() ([]models.SessionModel, error) {
	return nil, nil
}

func (sessionService *SessionService) CreateSession(session models.SessionModel) (*models.SessionModel, error) {
	result := sessionService.db.Create(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (sessionService *SessionService) GetByID(sessionId string) (*models.SessionModel, error) {
	var session models.SessionModel
	result := sessionService.db.Where("session_id = ?", sessionId).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (sessionService *SessionService) DeleteSession(sessionId string) error {
	result := sessionService.db.Where("session_id = ?", sessionId).Delete(&models.SessionModel{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sessionService *SessionService) InvalidateSession(c echo.Context, sessionID string) {
	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
	})

	err := sessionService.DeleteSession(sessionID)
	if err != nil {
		c.Logger().Error("Failed to delete session: ", err)
	}
}
