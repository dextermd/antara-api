package services

import (
	"antara-api/cmd/api/dtos"
	"antara-api/common"
	"antara-api/internal/models"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (sessionService *SessionService) CreateSession(userID uint, device, userAgent, ip string) (*models.SessionModel, error) {
	sessionID, err := common.GenerateSessionID()
	if err != nil {
		return nil, err
	}

	session := models.SessionModel{
		ID:        sessionID,
		UserID:    userID,
		Device:    device,
		UserAgent: userAgent,
		IPAddress: ip,
		IsActive:  true,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 дней
	}

	result := sessionService.db.Create(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

func (sessionService *SessionService) ValidateSession(sessionID string) (*models.UserModel, error) {
	var session models.SessionModel

	result := sessionService.db.Where("id = ? AND is_active = ? AND expires_at > ?", sessionID, true, time.Now()).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	var user models.UserModel
	result = sessionService.db.Preload("Roles").First(&user, session.UserID)
	if result.Error != nil {
		return nil, result.Error
	}

	if !user.IsActive {
		return nil, fmt.Errorf("аккаунт заблокирован")
	}

	sessionService.db.Model(&session).Update("updated_at", time.Now())
	return &user, nil
}

func (sessionService *SessionService) GetSessions(userID uint, currentSessionID string) ([]dtos.SessionInfo, error) {
	var sessions []models.SessionModel

	result := sessionService.db.Where("user_id = ? AND is_active = ?", userID, true).Order("created_at desc").Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}

	var currentSessionInfo []dtos.SessionInfo
	for _, session := range sessions {
		currentSessionInfo = append(currentSessionInfo, dtos.SessionInfo{
			ID:        session.ID,
			Device:    session.Device,
			UserAgent: session.UserAgent,
			IP:        session.IPAddress,
			IsActive:  session.IsActive,
			IsCurrent: session.ID == currentSessionID,
			CreatedAt: session.CreatedAt,
			LastUsed:  session.LastActivity,
		})
	}

	return currentSessionInfo, nil

}

func (sessionService *SessionService) GetUserFromSession(sessionID string) (*models.UserModel, error) {
	if sessionID == "" {
		return nil, nil
	}

	user, err := sessionService.ValidateSession(sessionID)
	if err != nil {
		return nil, nil
	}

	return user, nil
}

func (sessionService *SessionService) UpdateSessionActivity(sessionID string, userID uint) error {
	result := sessionService.db.Model(&models.SessionModel{}).Where("id = ? AND user_id = ?", sessionID, userID).Update("last_activity", time.Now())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sessionService *SessionService) DeleteSession(sessionID string, userID uint) error {
	result := sessionService.db.Where("id = ? AND user_id = ?", sessionID, userID).Delete(&models.SessionModel{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sessionService *SessionService) RevokeSession(sessionID string) error {
	result := sessionService.db.Where("id = ?", sessionID).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sessionService *SessionService) RevokeAllSessions(userID uint, exceptSessionID string) error {
	result := sessionService.db.Where("user_id = ? AND id != ?", userID, exceptSessionID).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
