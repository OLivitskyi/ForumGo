package sqlite

import (
	"Forum/internal/model"
	"fmt"
)

type ReactionRepository struct {
	store *Store
}

func (r *ReactionRepository) CreateReaction(reaction *model.Reaction) error {
	queryInsert := "INSERT INTO reactions(user_UUID, post_id, type) VALUES (?, ?, ?)"
	_, err := r.store.Db.Exec(queryInsert, reaction.UserID, reaction.PostID, reaction.Type)
	if err != nil {
		return fmt.Errorf("CreateReaction error: %w", err)
	}
	return nil
}

func (r *ReactionRepository) DeleteReaction(userID, postID string) error {
	// Create a query to delete the reaction
	queryDelete := "DELETE FROM reactions WHERE user_UUID = ? AND post_id = ?"
	if _, err := r.store.Db.Exec(queryDelete, userID, postID); err != nil {
		return fmt.Errorf("DeleteReaction error: %w", err)
	}

	return nil
}

func (r *ReactionRepository) GetUserReaction(userID, postID string) (*model.Reaction, error) {
	// Create a query to get the reaction
	var reaction model.Reaction
	queryGet := "SELECT user_UUID, post_id, type FROM reactions WHERE user_UUID = ? AND post_id = ?"
	if err := r.store.Db.QueryRow(queryGet, userID, postID).Scan(&reaction.UserID, &reaction.PostID, &reaction.Type); err != nil {
		return nil, fmt.Errorf("GetUserReaction error: %w", err)
	}

	return &reaction, nil
}

func (r *ReactionRepository) CountReactions(postID string) (int, error) {
	// Create a query to count the reactions (disregarding the type)
	queryCount := "SELECT COUNT(*) FROM reactions WHERE post_id = ?"
	var count int
	if err := r.store.Db.QueryRow(queryCount, postID).Scan(&count); err != nil {
		return 0, fmt.Errorf("CountReactions error: %w", err)
	}

	return count, nil
}
