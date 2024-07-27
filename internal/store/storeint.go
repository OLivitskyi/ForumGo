// internal/store/interfaces.go
package store

import (
	"Forum/internal/model"
	"time"
)

type Store interface {
	User() UserRepository
	Post() PostRepository
	Category() CategoryRepository
	Session() SessionRepository
	Comment() CommentRepository
	Reaction() ReactionRepository
	Message() MessageRepository
	UserStatus() UserStatusRepository
}

type CategoryRepository interface {
	Create(cate *model.Category) error
	GetAll() ([]*model.Category, error)
	AddCategoryToPost(postID string, categoryID int) error
	Exists(name string) (bool, error)
}

type SessionRepository interface {
	Create(s *model.Session) error
	GetByUUID(uuid string) (*model.Session, error)
	Delete(uuid string) error
	SaveSession(token string, userID int, expiration time.Time) error
	DeleteSession(token string) error
	GetUserIDFromSession(token string) (int, error)
}

type UserRepository interface {
	ExistingUser(userName, email string) error
	Login(user *model.User) error
	Register(user *model.User) error
	GetByUUID(uuid string) (*model.User, error)
	GetUsersOrderedByLastMessageOrAlphabetically() ([]*model.User, error)
}

type PostRepository interface {
	Create(post *model.Post) error
	GetAll() ([]*model.Post, error)
	AddCategoryToPost(postID string, categoryID int) error
	GetCategories(postID string) ([]*model.Category, error)
	GetByCategory(categoryID int) ([]*model.Post, error)
}

type CommentRepository interface {
	Create(c *model.Comment) error
	GetByPostID(postID string) ([]*model.Comment, error)
	GetCommentsWithReactionsByPostID(postID string) ([]*model.Comment, error)
}

type ReactionRepository interface {
	CreatePostReaction(reaction *model.Reaction) error
	DeletePostReaction(userID, postID string) error
	GetUserPostReaction(userID, postID string) (*model.Reaction, error)
	CountPostReactions(postID string) (int, error)
	UpdatePostReaction(userID, postID string, reactionID int) error

	CreateCommentReaction(reaction *model.Reaction) error
	DeleteCommentReaction(userID, commentID string) error
	GetUserCommentReaction(userID, commentID string) (*model.Reaction, error)
	CountCommentReactions(commentID string) (int, error)
}

type MessageRepository interface {
	AddMessage(senderID, receiverID int, content string) error
	GetMessages(senderID, receiverID int, limit, offset int) ([]*model.Message, error)
	MarkMessageAsRead(messageID int, userID int) error
}

type UserStatusRepository interface {
	UpdateUserStatus(userID int, isOnline bool) error
	GetUserStatus() ([]*model.UserStatus, error)
}
