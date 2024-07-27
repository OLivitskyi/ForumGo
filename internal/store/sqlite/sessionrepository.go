// internal/store/sqlite/sessionrepository.go
package sqlite

import (
	"Forum/internal/model"
	"time"
)

type SessionRepository struct {
	store *sqliteStore
}

func (r *SessionRepository) Create(s *model.Session) error {
	statement := "INSERT INTO sessions(user_UUID, session_id, expires_at) VALUES (?, ?, ?)"
	_, err := r.store.db.Exec(statement, s.UserUUID, s.SessionID, s.ExpiresAt)
	return err
}

func (r *SessionRepository) GetByUUID(sessionID string) (*model.Session, error) {
	var s model.Session
	if err := r.store.db.QueryRow(
		"SELECT user_UUID, session_id, expires_at FROM sessions WHERE session_id = ?",
		sessionID,
	).Scan(&s.UserUUID, &s.SessionID, &s.ExpiresAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) Delete(sessionID string) error {
	_, err := r.store.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

func (r *SessionRepository) SaveSession(token string, userID int, expiration time.Time) error {
	stmt, err := r.store.db.Prepare(`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token, userID, expiration.Unix())
	return err
}

func (r *SessionRepository) DeleteSession(token string) error {
	stmt, err := r.store.db.Prepare(`DELETE FROM sessions WHERE token = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token)
	return err
}

func (r *SessionRepository) GetUserIDFromSession(token string) (int, error) {
	var userID int
	err := r.store.db.QueryRow(`SELECT user_id FROM sessions WHERE token = ?`, token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
