package sqlite

import (
	"Forum/internal/model"
)

type MessageRepository struct {
	store *sqliteStore
}

func (r *MessageRepository) AddMessage(senderID, receiverID int, content string) error {
	stmt, err := r.store.db.Prepare(`INSERT INTO messages (sender_id, receiver_id, content, is_read) VALUES (?, ?, ?, 0)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(senderID, receiverID, content)
	return err
}

func (r *MessageRepository) GetMessages(senderID, receiverID int, limit, offset int) ([]*model.Message, error) {
	rows, err := r.store.db.Query(`SELECT message_id, sender_id, receiver_id, content, created_at, is_read FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY created_at DESC LIMIT ? OFFSET ?`, senderID, receiverID, receiverID, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &msg.IsRead)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (r *MessageRepository) MarkMessageAsRead(messageID int, userID int) error {
	stmt, err := r.store.db.Prepare(`UPDATE messages SET is_read = 1 WHERE message_id = ? AND receiver_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(messageID, userID)
	return err
}
