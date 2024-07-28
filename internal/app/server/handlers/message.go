package handlers

import (
	"Forum/internal/store"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// SendMessageHandler обробляє відправлення повідомлень
func SendMessageHandler(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		senderID, err := getUserIDFromSession(r, store)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		receiverID, err := strconv.Atoi(r.FormValue("receiver_id"))
		if err != nil {
			http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")
		if err := store.Message().AddMessage(senderID, receiverID, content); err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// GetMessagesHandler обробляє отримання повідомлень
func GetMessagesHandler(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromSession(r, store)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		otherUserID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 10
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			offset = 0
		}
		messages, err := store.Message().GetMessages(userID, otherUserID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to get messages", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

// MarkMessageAsReadHandler обробляє позначення повідомлень як прочитаних
func MarkMessageAsReadHandler(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromSession(r, store)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		messageID, err := strconv.Atoi(r.FormValue("message_id"))
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}
		if err := store.Message().MarkMessageAsRead(messageID, userID); err != nil {
			http.Error(w, "Failed to mark message as read", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
