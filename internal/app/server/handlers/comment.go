package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"html/template"
	"log"
	"net/http"
)

// CreateComment creates a new comment
func CreateComment(store store.Store, templates *template.Template, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		session, err := store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		userUUID := session.UserUUID
		postID := r.FormValue("postID")
		commentTxt := r.FormValue("commentText")

		comment, err := model.NewComment(postID, userUUID, commentTxt)
		if err != nil {
			logger.Println("NewComment() error: ", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err = store.Comment().Create(comment); err != nil {
			logger.Println("CreateComment() error: ", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
