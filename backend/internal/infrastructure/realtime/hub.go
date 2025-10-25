package realtime

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypeProgress     MessageType = "progress"
	MessageTypeError        MessageType = "error"
	MessageTypeSuccess      MessageType = "success"
	MessageTypeNotification MessageType = "notification"
	MessageTypeData         MessageType = "data"
	MessageTypeAuth         MessageType = "auth"
	MessageTypePing         MessageType = "ping"
	MessageTypePong         MessageType = "pong"
	MessageTypeRealtimeData MessageType = "realtime_data"
)

type Message struct {
	Type      MessageType            `json:"type"`
	ID        string                 `json:"id,omitempty"`
	Data      interface{}            `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ProgressData struct {
	OperationID string  `json:"operation_id"`
	Step        string  `json:"step"`
	Progress    float64 `json:"progress"`
	Message     string  `json:"message"`
	ETA         *int64  `json:"eta,omitempty"`
}

type Client struct {
	ID            string
	UserID        uuid.UUID
	Conn          WebSocketConnection
	Send          chan Message
	Hub           *Hub
	Context       context.Context
	Cancel        context.CancelFunc
	Authenticated bool
	LastPong      time.Time
	Subscriptions map[string]bool
	mu            sync.RWMutex
}

type WebSocketConnection struct {
	*websocket.Conn
}

func (w *WebSocketConnection) Close() error {
	return w.Conn.Close()
}

func (w *WebSocketConnection) SetReadDeadline(t time.Time) error {
	return w.Conn.SetReadDeadline(t)
}

func (w *WebSocketConnection) SetWriteDeadline(t time.Time) error {
	return w.Conn.SetWriteDeadline(t)
}

func (w *WebSocketConnection) ReadMessage() (messageType int, p []byte, err error) {
	return w.Conn.ReadMessage()
}

func (w *WebSocketConnection) WriteMessage(messageType int, data []byte) error {
	return w.Conn.WriteMessage(messageType, data)
}

type Hub struct {
	clients map[*Client]bool

	clientsByUser map[uuid.UUID]map[*Client]bool

	clientsByChannel map[string]map[*Client]bool

	register chan *Client

	unregister chan *Client

	broadcast chan Message

	sendToUser chan UserMessage

	sendToClient chan ClientMessage

	sendToChannel chan ChannelMessage

	mutex sync.RWMutex
}

type UserMessage struct {
	UserID  uuid.UUID
	Message Message
}

type ClientMessage struct {
	ClientID string
	Message  Message
}

type ChannelMessage struct {
	Channel string
	Message Message
}

func NewHub() *Hub {
	return &Hub{
		clients:          make(map[*Client]bool),
		clientsByUser:    make(map[uuid.UUID]map[*Client]bool),
		clientsByChannel: make(map[string]map[*Client]bool),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		broadcast:        make(chan Message),
		sendToUser:       make(chan UserMessage),
		sendToClient:     make(chan ClientMessage),
		sendToChannel:    make(chan ChannelMessage),
	}
}

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

		case channelMsg := <-h.sendToChannel:
			h.sendToChannelClients(channelMsg.Channel, channelMsg.Message)
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastMessage(message Message) {
	h.broadcast <- message
}

func (h *Hub) SendToUser(userID uuid.UUID, message Message) {
	h.sendToUser <- UserMessage{UserID: userID, Message: message}
}

func (h *Hub) SendToClient(clientID string, message Message) {
	h.sendToClient <- ClientMessage{ClientID: clientID, Message: message}
}

func (h *Hub) SendToChannel(channel string, message Message) {
	h.sendToChannel <- ChannelMessage{Channel: channel, Message: message}
}

func (h *Hub) SubscribeClient(clientID, channel string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		if client.ID == clientID {
			client.mu.Lock()
			if client.Subscriptions == nil {
				client.Subscriptions = make(map[string]bool)
			}
			client.Subscriptions[channel] = true
			client.mu.Unlock()

			if h.clientsByChannel[channel] == nil {
				h.clientsByChannel[channel] = make(map[*Client]bool)
			}
			h.clientsByChannel[channel][client] = true

			log.Printf("Client %s subscribed to channel: %s", clientID, channel)
			break
		}
	}
}

func (h *Hub) UnsubscribeClient(clientID, channel string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		if client.ID == clientID {
			client.mu.Lock()
			delete(client.Subscriptions, channel)
			client.mu.Unlock()

			if channelClients, exists := h.clientsByChannel[channel]; exists {
				delete(channelClients, client)
				if len(channelClients) == 0 {
					delete(h.clientsByChannel, channel)
				}
			}

			log.Printf("Client %s unsubscribed from channel: %s", clientID, channel)
			break
		}
	}
}

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

		client.mu.RLock()
		for channel := range client.Subscriptions {
			if channelClients, exists := h.clientsByChannel[channel]; exists {
				delete(channelClients, client)
				if len(channelClients) == 0 {
					delete(h.clientsByChannel, channel)
				}
			}
		}
		client.mu.RUnlock()

		log.Printf("Client %s unregistered", client.ID)
	}
}

func (h *Hub) broadcastMessage(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:

		}
	}
}

func (h *Hub) sendToUserClients(userID uuid.UUID, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if userClients, exists := h.clientsByUser[userID]; exists {
		for client := range userClients {
			select {
			case client.Send <- message:
			default:

			}
		}
	}
}

func (h *Hub) sendToSpecificClient(clientID string, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.ID == clientID {
			select {
			case client.Send <- message:
			default:

			}
			break
		}
	}
}

func (h *Hub) sendToChannelClients(channel string, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if channelClients, exists := h.clientsByChannel[channel]; exists {
		for client := range channelClients {
			select {
			case client.Send <- message:
			default:

			}
		}
	}
}

func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

func (h *Hub) GetUserClientCount(userID uuid.UUID) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if userClients, exists := h.clientsByUser[userID]; exists {
		return len(userClients)
	}
	return 0
}

func (h *Hub) GetChannelClientCount(channel string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if channelClients, exists := h.clientsByChannel[channel]; exists {
		return len(channelClients)
	}
	return 0
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
	}()

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

			var msg Message
			if err := json.Unmarshal(message, &msg); err == nil {
				if msg.Type == MessageTypePong {
					c.LastPong = time.Now()
				}
			}
		}
	}
}

func (c *Client) WritePump() {
	pingTicker := time.NewTicker(30 * time.Second)
	healthTicker := time.NewTicker(5 * time.Second)
	defer func() {
		pingTicker.Stop()
		healthTicker.Stop()
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

		case <-pingTicker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			ping := Message{Type: MessageTypePing, Timestamp: time.Now()}
			if err := c.Conn.WriteMessage(websocket.TextMessage, mustJSON(ping)); err != nil {
				return
			}

		case <-healthTicker.C:

			if time.Since(c.LastPong) > 60*time.Second {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

		case <-c.Context.Done():
			return
		}
	}
}

func mustJSON(m Message) []byte {
	b, err := json.Marshal(m)
	if err != nil {

		return []byte(`{"type":"ping"}`)
	}
	return b
}

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

	}
}

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

	}
}

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

	}
}
