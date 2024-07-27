package sqlite

import "Forum/internal/model"

type CommentRepository struct {
	store *sqliteStore
}

func (r *CommentRepository) Create(c *model.Comment) error {
	_, err := r.store.db.Exec(`
        INSERT INTO comments (id, post_id, user_UUID, content, created_at)
        VALUES (?, ?, ?, ?, datetime('now'))
    `, c.ID, c.PostID, c.UserID, c.Content)

	return err
}

func (r *CommentRepository) GetCommentsWithReactionsByPostID(postID string) ([]*model.Comment, error) {
	query := `
    SELECT c.id, c.post_id, c.user_uuid, c.content, c.created_at,
           COALESCE(SUM(CASE WHEN cr.reaction_id = 1 THEN 1 ELSE 0 END), 0) AS LikeCount,
           COALESCE(SUM(CASE WHEN cr.reaction_id = 2 THEN 1 ELSE 0 END), 0) AS DislikeCount
	FROM comments c
	LEFT JOIN comment_reactions cr ON c.id = cr.comment_id
	WHERE c.post_id = ?
	GROUP BY c.id`
	rows, err := r.store.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.LikeCount, &comment.DislikeCount)
		if err != nil {
			return nil, err
		}
		user, err := r.store.User().GetByUUID(comment.UserID)
		if err != nil {
			return nil, err
		}
		comment.User = user

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetByPostID(postID string) ([]*model.Comment, error) {
	rows, err := r.store.db.Query(`
		SELECT id, post_id, user_UUID, content, created_at 
		FROM comments
		WHERE post_id = ?
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*model.Comment, 0)
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}

		user, err := r.store.User().GetByUUID(c.UserID)
		if err != nil {
			return nil, err
		}
		c.User = user

		comments = append(comments, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
