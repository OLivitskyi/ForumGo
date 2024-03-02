package sqlite

import (
	"Forum/internal/model"
	"log"
)

type PostRepository struct {
	store  *Store
	Logger *log.Logger
}

func (r *PostRepository) AddCategoryToPost(postID string, categoryID int) error {
	_, err := r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
	return err
}

func (r *PostRepository) GetAll() ([]*model.Post, error) {
	rows, err := r.store.Db.Query("SELECT id, user_UUID, subject, content FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*model.Post, 0)
	for rows.Next() {
		var p model.Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Subject, &p.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetCategories(postID string) ([]*model.Category, error) {
	rows, err := r.store.Db.Query(`
        SELECT categories.id, categories.category_name
        FROM categories, post_categories
        WHERE post_categories.post_id = ?
        AND post_categories.category_id = categories.id
    `, postID)
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
	return categories, rows.Err()
}

func (r *PostRepository) GetByCategory(categoryID int) ([]*model.Post, error) {
	rows, err := r.store.Db.Query(`
    SELECT posts.id, posts.user_UUID, posts.subject, posts.content
    FROM posts
    INNER JOIN post_categories ON posts.id = post_categories.post_id
    WHERE post_categories.category_id = ?
`, categoryID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := make([]*model.Post, 0)
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Subject, &post.Content)
		if err != nil {
			return nil, err
		}

		categories, err := r.GetCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, &post)
	}

	return posts, rows.Err()
}
