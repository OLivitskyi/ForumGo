package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// CreatePost створює новий пост
func CreatePost(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримання cookie сесії
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Отримання сесії
		session, err := store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Отримання UUID користувача з сесії
		userUUID := session.UserUUID

		subject := r.FormValue("postTitle")
		content := r.FormValue("postText")

		// Розбір форми для чекбоксів категорій
		r.ParseForm()
		categoryIDs := r.PostForm["categoryIDs"]

		// Перевірка, чи вибрана хоча б одна категорія
		if len(categoryIDs) == 0 {
			http.Redirect(w, r, "/createPostPage?error=atleast_one_category_required", http.StatusSeeOther)
			return
		}

		post, err := model.NewPost(userUUID, subject, content)
		if err != nil {
			logger.Println("NewPost() error: ", err)
			http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
			return
		}

		if err = store.Post().Create(post); err != nil {
			logger.Println("Create() error: ", err)
			http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
			return
		}

		// Додавання категорій до поста
		for _, categoryIDStr := range categoryIDs {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				logger.Println("Error converting categoryID to int: ", err)
				http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
				return
			}

			if err := store.Post().AddCategoryToPost(post.ID, categoryID); err != nil {
				logger.Println("Error adding category to post: ", err)
				http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// CreatePostPage рендерить сторінку створення поста
func CreatePostPage(store store.Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := store.Category().GetAll()
		if err != nil {
			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return
		}

		errMsg := ""
		if errParam := r.URL.Query().Get("error"); errParam == "atleast_one_category_required" {
			errMsg = "At least one category must be selected."
		}

		data := struct {
			Categories   []*model.Category
			ErrorMessage string
		}{
			Categories:   categories,
			ErrorMessage: errMsg,
		}

		if err := templates.Lookup("createPostPage.html").Execute(w, data); err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		}
	}
}
