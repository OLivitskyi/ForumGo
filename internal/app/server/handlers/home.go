package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"html/template"
	"log"
	"net/http"
)

func Home(store store.Store, templates *template.Template, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current user if exists
		var user *model.User
		if sessionCookie, err := r.Cookie("session_uuid"); err == nil {
			session, err := store.Session().GetByUUID(sessionCookie.Value)
			if err == nil {
				user, _ = store.User().GetByUUID(session.UserUUID)
			}
		}

		// Fetching all posts
		posts, err := store.Post().GetAll()
		if err != nil {
			logger.Println("error fetching posts:", err)
			http.Error(w, "error fetching posts", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			// Fetch user who created the post
			fetchedUser, _ := store.User().GetByUUID(post.UserID)
			post.User = fetchedUser

			// Fetch categories for each post
			categories, err := store.Post().GetCategories(post.ID)
			if err != nil {
				logger.Println("error fetching categories for post:", err)
				http.Error(w, "error fetching post categories", http.StatusInternalServerError)
				return
			}
			post.Categories = categories

			// Fetch comments with reactions for each post using the updated repository method
			comments, err := store.Comment().GetCommentsWithReactionsByPostID(post.ID)
			if err != nil {
				logger.Println("error fetching comments for post:", err)
				http.Error(w, "error fetching post comments", http.StatusInternalServerError)
				return
			}
			for _, comment := range comments {
				fetchedUser, _ := store.User().GetByUUID(comment.UserID)
				comment.User = fetchedUser
			}
			post.Comments = comments
		}

		// Fetching all categories
		allCategories, err := store.Category().GetAll()
		if err != nil {
			logger.Println("error fetching categories:", err)
			http.Error(w, "error fetching categories", http.StatusInternalServerError)
			return
		}

		// Struct to pass into template
		data := &model.PageData{
			User:       user,
			Posts:      posts,
			Categories: allCategories,
		}

		execTmpl(w, templates.Lookup("home.html"), data)
	}
}
