package model

import "time"

type Session struct {
	ID        int
	UserUUID  string
	SessionID string
	ExpiresAt time.Time
}

// NewSession takes a user UUID and returns a new session struct
func NewSession(uuid string) (*Session, error) {
	// Calculate expiry time for the session. Here we set it to 1 hour ahead of current time.
	expiresAt := time.Now().Add(time.Hour * 1)

	return &Session{
		UserUUID:  uuid,
		SessionID: uuid,
		ExpiresAt: expiresAt,
	}, nil
}
