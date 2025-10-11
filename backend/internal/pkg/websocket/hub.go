package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeProgress     MessageType = "progress"
	MessageTypeError        MessageType = "error"
	MessageTypeSuccess      MessageType = "success"
	MessageTypeNotification MessageType = "notification"
	MessageTypeData         MessageType = "data"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType            `json:"type"`
	ID        string                 `json:"id,omitempty"`
	Data      interface{}            `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ProgressData represents progress information for long-running operations
type ProgressData struct {
	OperationID string  `json:"operation_id"`
	Step        string  `json:"step"`
	Progress    float64 `json:"progress"` // 0.0 to 1.0
	Message     string  `json:"message"`
	ETA         *int64  `json:"eta,omitempty"` // seconds until completion
}

// Client represents a WebSocket client connection
type Client struct {
	ID      string
	UserID  uuid.UUID
	Conn    WebSocketConnection
	Send    chan Message
	Hub     *Hub
	Context context.Context
	Cancel  context.CancelFunc
}

// WebSocketConnection wraps the gorilla websocket connection
type WebSocketConnection struct {
	*websocket.Conn
}

// Close closes the WebSocket connection
func (w *WebSocketConnection) Close() error {
	return w.Conn.Close()
}

// SetReadDeadline sets the read deadline
func (w *WebSocketConnection) SetReadDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline
func (w *WebSocketConnection) SetWriteDeadline(t time.Time) error {
	return w.Conn.SetWriteDeadline(t)
}

// ReadMessage reads a message from the connection
func (w *WebSocketConnection) ReadMessage() (messageType int, p []byte, err error) {
	return w.Conn.ReadMessage()
}

// WriteMessage writes a message to the connection
func (w *WebSocketConnection) WriteMessage(messageType int, data []byte) error {
	return w.Conn.WriteMessage(messageType, data)
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Clients by user ID
	clientsByUser map[uuid.UUID]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to all clients
	broadcast chan Message

	// Send message to specific user
	sendToUser chan UserMessage

	// Send message to specific client
	sendToClient chan ClientMessage

	// Mutex for thread-safe operations
	mutex sync.RWMutex
}

// UserMessage represents a message to be sent to a specific user
type UserMessage struct {
	UserID  uuid.UUID
	Message Message
}

// ClientMessage represents a message to be sent to a specific client
type ClientMessage struct {
	ClientID string
	Message  Message
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		clientsByUser: make(map[uuid.UUID]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan Message),
		sendToUser:    make(chan UserMessage),
		sendToClient:  make(chan ClientMessage),
	}
}

// Run starts the hub
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case userMsg := <-h.sendToUser:
			h.sendToUserClients(userMsg.UserID, userMsg.Message)

		case clientMsg := <-h.sendToClient:
			h.sendToSpecificClient(clientMsg.ClientID, clientMsg.Message)
		}
	}
}

// RegisterClient registers a new client
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastMessage broadcasts a message to all clients
func (h *Hub) BroadcastMessage(message Message) {
	h.broadcast <- message
}

// SendToUser sends a message to all clients of a specific user
func (h *Hub) SendToUser(userID uuid.UUID, message Message) {
	h.sendToUser <- UserMessage{UserID: userID, Message: message}
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(clientID string, message Message) {
	h.sendToClient <- ClientMessage{ClientID: clientID, Message: message}
}

// registerClient handles client registration
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true

	if h.clientsByUser[client.UserID] == nil {
		h.clientsByUser[client.UserID] = make(map[*Client]bool)
	}
	h.clientsByUser[client.UserID][client] = true

	log.Printf("Client %s registered for user %s", client.ID, client.UserID)
}

// unregisterClient handles client unregistration
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		if userClients, exists := h.clientsByUser[client.UserID]; exists {
			delete(userClients, client)
			if len(userClients) == 0 {
				delete(h.clientsByUser, client.UserID)
			}
		}

		log.Printf("Client %s unregistered", client.ID)
	}
}

// broadcastMessage broadcasts a message to all clients
func (h *Hub) broadcastMessage(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// sendToUserClients sends a message to all clients of a specific user
func (h *Hub) sendToUserClients(userID uuid.UUID, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if userClients, exists := h.clientsByUser[userID]; exists {
		for client := range userClients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
				delete(userClients, client)
			}
		}
	}
}

// sendToSpecificClient sends a message to a specific client
func (h *Hub) sendToSpecificClient(clientID string, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.ID == clientID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
			break
		}
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetUserClientCount returns the number of clients for a specific user
func (h *Hub) GetUserClientCount(userID uuid.UUID) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if userClients, exists := h.clientsByUser[userID]; exists {
		return len(userClients)
	}
	return 0
}

// Client methods

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	for {
		select {
		case <-c.Context.Done():
			return
		default:
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			// Handle incoming message (ping/pong, etc.)
			var msg Message
			if err := json.Unmarshal(message, &msg); err == nil {
				// Echo back ping messages as pong
				if msg.Type == "ping" {
					pongMsg := Message{
						Type:      "pong",
						ID:        msg.ID,
						Timestamp: time.Now(),
					}
					c.Send <- pongMsg
				}
			}
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.Context.Done():
			return
		}
	}
}

// SendProgress sends a progress update to the client
func (c *Client) SendProgress(operationID, step string, progress float64, message string) {
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

	select {
	case c.Send <- msg:
	default:
		// Client is not ready to receive messages
	}
}

// SendError sends an error message to the client
func (c *Client) SendError(operationID, message string, err error) {
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

	select {
	case c.Send <- msg:
	default:
		// Client is not ready to receive messages
	}
}

// SendSuccess sends a success message to the client
func (c *Client) SendSuccess(operationID, message string, data interface{}) {
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

	select {
	case c.Send <- msg:
	default:
		// Client is not ready to receive messages
	}
}
