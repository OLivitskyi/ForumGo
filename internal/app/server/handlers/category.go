package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// CreateCategoryPage renders the create category page
func CreateCategoryPage(store store.Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := ""
		switch r.URL.Query().Get("error") {
		case "checkError":
			errorMessage = "Failed to check category existence."
		case "categoryExists":
			errorMessage = "This category already exists."
		}
		execTmpl(w, templates.Lookup("createCategory.html"), map[string]string{"ErrorMessage": errorMessage})
	}
}

// CreateCategory handles the creation of a new category
func CreateCategory(store store.Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryName := r.FormValue("categoryName")
		exists, err := store.Category().Exists(categoryName)
		if err != nil {
			http.Redirect(w, r, "/createCategoryPage?error=checkError", http.StatusSeeOther)
			return
		}

		if exists {
			http.Redirect(w, r, "/createCategoryPage?error=categoryExists", http.StatusSeeOther)
			return
		}

		category := &model.Category{Name: categoryName}
		if err := store.Category().Create(category); err != nil {
			http.Redirect(w, r, "/createCategoryPage", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// CategoryPosts renders posts in a specific category
func CategoryPosts(store store.Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *model.User
		if sessionCookie, err := r.Cookie("session_uuid"); err == nil {
			session, err := store.Session().GetByUUID(sessionCookie.Value)
			if err == nil {
				user, _ = store.User().GetByUUID(session.UserUUID)
			}
		}

		categoryIDStr := strings.TrimPrefix(r.URL.Path, "/category/")
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		posts, err := store.Post().GetByCategory(categoryID)
		if err != nil {
			http.Error(w, "error fetching posts by category", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			fetchedUser, _ := store.User().GetByUUID(post.UserID)
			post.User = fetchedUser

			categories, err := store.Post().GetCategories(post.ID)
			if err != nil {
				http.Error(w, "error fetching post categories", http.StatusInternalServerError)
				return
			}
			post.Categories = categories

			commentsWithReactions, err := store.Comment().GetCommentsWithReactionsByPostID(post.ID)
			if err != nil {
				http.Error(w, "error fetching post comments with reactions", http.StatusInternalServerError)
				return
			}
			for _, comment := range commentsWithReactions {
				fetchedUser, _ := store.User().GetByUUID(comment.UserID)
				comment.User = fetchedUser
			}
			post.Comments = commentsWithReactions
		}

		allCategories, err := store.Category().GetAll()
		if err != nil {
			http.Error(w, "error fetching all categories", http.StatusInternalServerError)
			return
		}

		data := &model.PageData{
			User:       user,
			Posts:      posts,
			Categories: allCategories,
		}

		execTmpl(w, templates.Lookup("home.html"), data)
	}
}
