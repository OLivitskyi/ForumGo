package store

import "Forum/internal/model"

type Store interface {
	User() UserRepository
	Post() PostRepository
	Category() CategoryRepository
}
type CategoryRepository interface {
	Create(cate *model.Category) error
	GetAll() ([]*model.Category, error)
	AddCategoryToPost(postID string, categoryID int) error
}
