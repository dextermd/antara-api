package response

import (
	"antara-api/internal/models"
	"time"
)

type AuthDataResponse struct {
	User      *models.UserModel `json:"user"`
	SessionID string            `json:"session_id"`
	ExpiresAt time.Time         `json:"expires_at"`
}
