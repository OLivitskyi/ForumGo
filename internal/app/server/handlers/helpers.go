package handlers

import (
	"Forum/internal/store"
	"html/template"
	"log"
	"net/http"
)

// execTmpl is a helper function to execute templates
func execTmpl(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
	}
}

// getUserIDFromSession отримує ID користувача з сесії
func getUserIDFromSession(r *http.Request, store store.Store) (int, error) {
	cookie, err := r.Cookie("session_uuid")
	if err != nil {
		return 0, err
	}
	sessionToken := cookie.Value
	return store.Session().GetUserIDFromSession(sessionToken)
}
