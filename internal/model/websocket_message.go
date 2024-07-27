package model

// WebSocketMessage represents a message sent over WebSocket.
type WebSocketMessage struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Sender    int    `json:"sender"`
	Receiver  int    `json:"receiver"`
	Timestamp string `json:"timestamp"`
	IsRead    bool   `json:"is_read"`
}
