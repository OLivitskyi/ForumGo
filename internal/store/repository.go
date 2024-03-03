package store

import "Forum/internal/model"

type UserRepository interface {
	ExistingUser(userName, email string) error
	Login(user *model.User) error
	Register(user *model.User) error
	GetByUUID(uuid string) (*model.User, error)
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
}
