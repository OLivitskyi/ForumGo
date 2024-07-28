package handlers

import (
	"Forum/internal/store"
	"encoding/json"
	"log"
	"net/http"
)

// UpdateStatusHandler обробляє оновлення статусу користувача
func UpdateStatusHandler(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromSession(r, store)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		isOnline := r.FormValue("is_online") == "true"
		if err := store.UserStatus().UpdateUserStatus(userID, isOnline); err != nil {
			http.Error(w, "Failed to update status", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// GetUserStatusHandler обробляє отримання статусу користувачів
func GetUserStatusHandler(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statuses, err := store.UserStatus().GetUserStatus()
		if err != nil {
			http.Error(w, "Failed to get user statuses", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)
	}
}
