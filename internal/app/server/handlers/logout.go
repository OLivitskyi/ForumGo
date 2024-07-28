package handlers

import (
	"Forum/internal/store"
	"log"
	"net/http"
)

// Logout обробляє вихід користувача
func Logout(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Видалення сесії з бази даних
		err = store.Session().Delete(session.Value)
		if err != nil {
			http.Error(w, "Failed to end session", http.StatusInternalServerError)
			return
		}

		// Видалення cookie сесії
		session.MaxAge = -1
		http.SetCookie(w, session)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
