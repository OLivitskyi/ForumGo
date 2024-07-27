package sqlite

import (
	"Forum/internal/store"
	"database/sql"
)

type sqliteStore struct {
	db                 *sql.DB
	userRepository     *UserRepository
	postRepository     *PostRepository
	categoryRepository *CategoryRepository
	sessionRepository  *SessionRepository
	commentRepository  *CommentRepository
	reactionRepo       *ReactionRepository
	messageRepository  *MessageRepository
	userStatusRepo     *UserStatusRepository
}

func NewSQL(db *sql.DB) *sqliteStore {
	return &sqliteStore{
		db: db,
	}
}

func (s *sqliteStore) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{store: s}
	}
	return s.userRepository
}

func (s *sqliteStore) Post() store.PostRepository {
	if s.postRepository == nil {
		s.postRepository = &PostRepository{store: s}
	}
	return s.postRepository
}

func (s *sqliteStore) Category() store.CategoryRepository {
	if s.categoryRepository == nil {
		s.categoryRepository = &CategoryRepository{store: s}
	}
	return s.categoryRepository
}

func (s *sqliteStore) Session() store.SessionRepository {
	if s.sessionRepository == nil {
		s.sessionRepository = &SessionRepository{store: s}
	}
	return s.sessionRepository
}

func (s *sqliteStore) Comment() store.CommentRepository {
	if s.commentRepository == nil {
		s.commentRepository = &CommentRepository{store: s}
	}
	return s.commentRepository
}

func (s *sqliteStore) Reaction() store.ReactionRepository {
	if s.reactionRepo == nil {
		s.reactionRepo = &ReactionRepository{store: s}
	}
	return s.reactionRepo
}

func (s *sqliteStore) Message() store.MessageRepository {
	if s.messageRepository == nil {
		s.messageRepository = &MessageRepository{store: s}
	}
	return s.messageRepository
}

func (s *sqliteStore) UserStatus() store.UserStatusRepository {
	if s.userStatusRepo == nil {
		s.userStatusRepo = &UserStatusRepository{store: s}
	}
	return s.userStatusRepo
}
