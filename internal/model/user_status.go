package model

import "time"

// UserStatus represents the online status of a user.
type UserStatus struct {
	UserID       int       `json:"user_id"`
	IsOnline     bool      `json:"is_online"`
	LastActivity time.Time `json:"last_activity"`
}
