package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"database/sql"
	"log"
	"net/http"
)

// HandleCreatePostReaction обробляє створення реакцій на пости
func HandleCreatePostReaction(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Error(w, "Session error", http.StatusBadRequest)
			return
		}
		session, err := store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Session retrieval error", http.StatusInternalServerError)
			return
		}
		userUUID := session.UserUUID
		postID := r.FormValue("postID")
		reactionTypeStr := r.FormValue("reactionType")

		var reactionID int
		switch reactionTypeStr {
		case "like":
			reactionID = 1
		case "dislike":
			reactionID = 2
		default:
			http.Error(w, "Invalid reaction type", http.StatusBadRequest)
			return
		}

		existingReaction, err := store.Reaction().GetUserPostReaction(userUUID, postID)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if existingReaction != nil {
			if existingReaction.ReactionID == reactionID {
				if err := store.Reaction().DeletePostReaction(userUUID, postID); err != nil {
					http.Error(w, "Failed to delete reaction", http.StatusInternalServerError)
					return
				}
			} else {
				if err := store.Reaction().DeletePostReaction(userUUID, postID); err != nil {
					http.Error(w, "Failed to delete existing reaction", http.StatusInternalServerError)
					return
				}

				reaction := &model.Reaction{
					UserID:     userUUID,
					PostID:     postID,
					ReactionID: reactionID,
				}
				if err := store.Reaction().CreatePostReaction(reaction); err != nil {
					http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
					return
				}
			}
		} else {
			reaction := &model.Reaction{
				UserID:     userUUID,
				PostID:     postID,
				ReactionID: reactionID,
			}
			if err := store.Reaction().CreatePostReaction(reaction); err != nil {
				http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/home?postID="+postID, http.StatusSeeOther)
	}
}

// HandleCreateCommentReaction обробляє створення реакцій на коментарі
func HandleCreateCommentReaction(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Error(w, "Session error", http.StatusBadRequest)
			return
		}

		session, err := store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Session retrieval error", http.StatusInternalServerError)
			return
		}

		userUUID := session.UserUUID
		commentID := r.FormValue("commentID")
		reactionTypeStr := r.FormValue("reactionType")

		var reactionID int
		switch reactionTypeStr {
		case "like":
			reactionID = 1
		case "dislike":
			reactionID = 2
		default:
			http.Error(w, "Invalid reaction type", http.StatusBadRequest)
			return
		}

		existingReaction, err := store.Reaction().GetUserCommentReaction(userUUID, commentID)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if existingReaction != nil {
			if existingReaction.ReactionID == reactionID {
				if err := store.Reaction().DeleteCommentReaction(userUUID, commentID); err != nil {
					http.Error(w, "Failed to delete reaction", http.StatusInternalServerError)
					return
				}
			} else {
				if err := store.Reaction().DeleteCommentReaction(userUUID, commentID); err != nil {
					http.Error(w, "Failed to delete existing reaction", http.StatusInternalServerError)
					return
				}

				reaction := &model.Reaction{
					UserID:     userUUID,
					CommentID:  commentID,
					ReactionID: reactionID,
				}
				if err := store.Reaction().CreateCommentReaction(reaction); err != nil {
					http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
					return
				}
			}
		} else {
			reaction := &model.Reaction{
				UserID:     userUUID,
				CommentID:  commentID,
				ReactionID: reactionID,
			}
			if err := store.Reaction().CreateCommentReaction(reaction); err != nil {
				http.Error(w, "Failed to create reaction", http.StatusInternalServerError)
				return
			}
		}

		// Перенаправлення на попередню сторінку або на маршрут за замовчуванням, якщо не доступно
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/home"
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
	}
}
