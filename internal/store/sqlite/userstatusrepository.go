package sqlite

import (
	"Forum/internal/model"
)

type UserStatusRepository struct {
	store *sqliteStore
}

func (r *UserStatusRepository) UpdateUserStatus(userID int, isOnline bool) error {
	stmt, err := r.store.db.Prepare(`INSERT INTO user_status (user_id, is_online) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET is_online = ?, last_activity = CURRENT_TIMESTAMP`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, isOnline, isOnline)
	return err
}

func (r *UserStatusRepository) GetUserStatus() ([]*model.UserStatus, error) {
	rows, err := r.store.db.Query(`SELECT user_id, is_online, last_activity FROM user_status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var statuses []*model.UserStatus
	for rows.Next() {
		var status model.UserStatus
		err := rows.Scan(&status.UserID, &status.IsOnline, &status.LastActivity)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, &status)
	}
	return statuses, nil
}
