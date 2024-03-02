package model

import "github.com/gofrs/uuid"

type Post struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Subject    string      `json:"subject"`
	Content    string      `json:"content"`
	Categories []*Category `json:"categories"` // Updated to []*Category
}

func NewPost(userID, subject, content string) (*Post, error) {
	// Create a new UUID for the post
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &Post{
		ID:         id.String(),
		UserID:     userID,
		Subject:    subject,
		Content:    content,
		Categories: make([]*Category, 0), // Initialize slice of *Category
	}, nil
}
