package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]int)
var broadcast = make(chan model.WebSocketMessage)
var mutex = &sync.Mutex{}

func HandleConnections(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("Error upgrading to websocket: %v", err)
			return
		}
		defer func() {
			ws.Close()
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
		}()

		sessionToken := r.URL.Query().Get("session_token")
		if sessionToken == "" {
			log.Println("Unauthorized access: session token missing")
			return
		}

		userID, err := store.Session().GetUserIDFromSession(sessionToken)
		if err != nil {
			log.Printf("Unauthorized access: %v", err)
			return
		}

		log.Printf("User %d connected", userID)
		mutex.Lock()
		clients[ws] = userID
		mutex.Unlock()
		store.UserStatus().UpdateUserStatus(userID, true)

		for {
			var msg model.WebSocketMessage
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Printf("Error reading websocket message: %v", err)
				break
			}
			msg.Timestamp = time.Now().Format(time.RFC3339)
			broadcast <- msg
		}
	}
}

func HandleMessages(store store.Store) {
	for {
		msg := <-broadcast
		if err := store.Message().AddMessage(msg.Sender, msg.Receiver, msg.Content); err != nil {
			log.Printf("Error storing message in the database: %v", err)
			continue
		}
		for client := range clients {
			if clients[client] == msg.Receiver {
				err := client.WriteJSON(msg)
				if err != nil {
					client.Close()
					mutex.Lock()
					delete(clients, client)
					mutex.Unlock()
				}
			}
		}
	}
}
