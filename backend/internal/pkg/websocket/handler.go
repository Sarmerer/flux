package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
			return
		}

		// Extract user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			conn.Close()
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Create client context
		ctx, cancel := context.WithCancel(r.Context())

		// Create client
		client := &Client{
			ID:      uuid.New().String(),
			UserID:  userID,
			Conn:    WebSocketConnection{Conn: conn},
			Send:    make(chan Message, 256),
			Hub:     hub,
			Context: ctx,
			Cancel:  cancel,
		}

		// Register client with hub
		hub.RegisterClient(client)

		// Start goroutines for reading and writing
		go client.WritePump()
		go client.ReadPump()

		// Send welcome message
		welcomeMsg := Message{
			Type:      MessageTypeNotification,
			Data:      map[string]interface{}{"message": "Connected to Flow WebSocket"},
			Timestamp: time.Now(),
		}
		client.Send <- welcomeMsg
	}
}

// SendProgressToUser sends progress updates to a specific user
func SendProgressToUser(hub *Hub, userID uuid.UUID, operationID, step string, progress float64, message string) {
	progressData := ProgressData{
		OperationID: operationID,
		Step:        step,
		Progress:    progress,
		Message:     message,
	}

	msg := Message{
		Type:      MessageTypeProgress,
		ID:        operationID,
		Data:      progressData,
		Timestamp: time.Now(),
	}

	hub.SendToUser(userID, msg)
}

// SendErrorToUser sends error messages to a specific user
func SendErrorToUser(hub *Hub, userID uuid.UUID, operationID, message string, err error) {
	errorData := map[string]interface{}{
		"operation_id": operationID,
		"message":      message,
		"error":        err.Error(),
	}

	msg := Message{
		Type:      MessageTypeError,
		ID:        operationID,
		Data:      errorData,
		Timestamp: time.Now(),
	}

	hub.SendToUser(userID, msg)
}

// SendSuccessToUser sends success messages to a specific user
func SendSuccessToUser(hub *Hub, userID uuid.UUID, operationID, message string, data interface{}) {
	successData := map[string]interface{}{
		"operation_id": operationID,
		"message":      message,
		"result":       data,
	}

	msg := Message{
		Type:      MessageTypeSuccess,
		ID:        operationID,
		Data:      successData,
		Timestamp: time.Now(),
	}

	hub.SendToUser(userID, msg)
}

// SendNotificationToUser sends notification messages to a specific user
func SendNotificationToUser(hub *Hub, userID uuid.UUID, message string, data interface{}) {
	notificationData := map[string]interface{}{
		"message": message,
		"data":    data,
	}

	msg := Message{
		Type:      MessageTypeNotification,
		Data:      notificationData,
		Timestamp: time.Now(),
	}

	hub.SendToUser(userID, msg)
}
