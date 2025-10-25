package realtime

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {

		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocket(hub *Hub, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())

		client := &Client{
			ID: uuid.New().String(),

			Conn:          WebSocketConnection{Conn: conn},
			Send:          make(chan Message, 256),
			Hub:           hub,
			Context:       ctx,
			Cancel:        cancel,
			Authenticated: false,
			LastPong:      time.Now(),
			Subscriptions: make(map[string]bool),
		}

		hub.RegisterClient(client)

		go client.WritePump()
		go client.ReadPump()

		authTimer := time.NewTimer(5 * time.Second)
		defer authTimer.Stop()

		authenticated := make(chan struct{})

		go func() {
			_, p, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var msg Message
			if err := json.Unmarshal(p, &msg); err != nil {
				return
			}
			if msg.Type != MessageTypeAuth {
				return
			}
			m, ok := msg.Data.(map[string]interface{})
			if !ok {
				return
			}
			tokenRaw, ok := m["token"].(string)
			if !ok {
				return
			}
			if !strings.HasPrefix(tokenRaw, "Bearer ") {
				return
			}
			tokenString := strings.TrimPrefix(tokenRaw, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return
			}
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				return
			}
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return
			}

			client.UserID = userID
			client.Authenticated = true
			close(authenticated)
		}()

		select {
		case <-authenticated:

			welcomeMsg := Message{
				Type:      MessageTypeNotification,
				Data:      map[string]interface{}{"message": "Authenticated WebSocket connected", "client_id": client.ID},
				Timestamp: time.Now(),
			}
			client.Send <- welcomeMsg
		case <-authTimer.C:

			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "auth required"), time.Now().Add(1*time.Second))
			conn.Close()
			return
		}
	}
}

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
