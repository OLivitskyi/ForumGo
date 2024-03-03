package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type Post struct {
	ID         string
	UserID     string
	User       *User
	Subject    string
	Content    string
	Categories []*Category
	Comments   []*Comment
	CreatedAt  time.Time
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
		CreatedAt:  time.Now(),           // Add current time
		Categories: make([]*Category, 0), // Initialize slice of *Category
	}, nil
}
