package sqlite

import (
	"Forum/internal/model"
)

type CategoryRepository struct {
	store *sqliteStore
}

func (r *CategoryRepository) Create(cate *model.Category) error {
	_, err := r.store.db.Exec(`INSERT INTO categories (category_name) VALUES (?)`, cate.Name)
	return err
}

func (r *CategoryRepository) GetAll() ([]*model.Category, error) {
	rows, err := r.store.db.Query("SELECT id, category_name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) AddCategoryToPost(postID string, categoryID int) error {
	_, err := r.store.db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
	return err
}

func (r *CategoryRepository) Exists(name string) (bool, error) {
	var count int
	err := r.store.db.QueryRow(`SELECT COUNT(*) FROM categories WHERE category_name = ?`, name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
