package dtos

import "time"

type SessionInfo struct {
	ID        string    `json:"id"`
	Device    string    `json:"device"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	IsActive  bool      `json:"is_active"`
	IsCurrent bool      `json:"is_current"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}
