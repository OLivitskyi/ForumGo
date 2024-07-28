package handlers

import (
	"Forum/internal/store"
	"net/http"
)

// getUserIDFromSession отримує ID користувача з сесії
func getUserIDFromSession(r *http.Request, store store.Store) (int, error) {
	cookie, err := r.Cookie("session_uuid")
	if err != nil {
		return 0, err
	}
	sessionToken := cookie.Value
	return store.Session().GetUserIDFromSession(sessionToken)
}
