package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"html/template"
	"log"
	"net/http"
)

// Home обробляє головну сторінку
func Home(store store.Store, tmpl *template.Template, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримання поточного користувача, якщо він існує
		var user *model.User
		if sessionCookie, err := r.Cookie("session_uuid"); err == nil {
			session, err := store.Session().GetByUUID(sessionCookie.Value)
			if err == nil {
				user, _ = store.User().GetByUUID(session.UserUUID)
			}
		}

		// Отримання всіх постів
		posts, err := store.Post().GetAll()
		if err != nil {
			logger.Println("error fetching posts:", err)
			http.Error(w, "error fetching posts", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			// Отримання користувача, який створив пост
			fetchedUser, _ := store.User().GetByUUID(post.UserID)
			post.User = fetchedUser

			// Отримання категорій для кожного поста
			categories, err := store.Post().GetCategories(post.ID)
			if err != nil {
				logger.Println("error fetching categories for post:", err)
				http.Error(w, "error fetching post categories", http.StatusInternalServerError)
				return
			}
			post.Categories = categories

			// Отримання коментарів з реакціями для кожного поста
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

		// Отримання всіх категорій
		allCategories, err := store.Category().GetAll()
		if err != nil {
			logger.Println("error fetching categories:", err)
			http.Error(w, "error fetching categories", http.StatusInternalServerError)
			return
		}

		// Передача даних до шаблону
		data := &model.PageData{
			User:       user,
			Posts:      posts,
			Categories: allCategories,
		}

		execTmpl(w, tmpl.Lookup("main.html"), data)
	}
}

// execTmpl рендерить шаблон з переданими даними або повертає внутрішню помилку сервера
func execTmpl(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
	}
}
