package realtime

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Notifier provides flexible WebSocket notifications
type Notifier struct {
	hub *Hub
}

// NewNotifier creates a notifier
func NewNotifier(hub *Hub) *Notifier {
	return &Notifier{hub: hub}
}

// Notify sends a message to a specific user
func (n *Notifier) Notify(userID uuid.UUID, messageType string, data interface{}) {
	if n.hub == nil {
		return // Gracefully handle nil hub
	}

	msg := Message{
		Type:      MessageType(messageType),
		Data:      data,
		Timestamp: time.Now(),
	}

	n.hub.SendToUser(userID, msg)
}

// NotifyWithID sends a message to a specific user with an operation ID
func (n *Notifier) NotifyWithID(userID uuid.UUID, messageType string, operationID string, data interface{}) {
	if n.hub == nil {
		return // Gracefully handle nil hub
	}

	msg := Message{
		Type:      MessageType(messageType),
		ID:        operationID,
		Data:      data,
		Timestamp: time.Now(),
	}

	n.hub.SendToUser(userID, msg)
}

// NotifyAll sends to all connected users
func (n *Notifier) NotifyAll(messageType string, data interface{}) {
	if n.hub == nil {
		return
	}

	msg := Message{
		Type:      MessageType(messageType),
		Data:      data,
		Timestamp: time.Now(),
	}

	n.hub.BroadcastMessage(msg)
}

// NotifyChannel sends a message to all clients subscribed to a channel
func (n *Notifier) NotifyChannel(channel string, messageType string, data interface{}) {
	if n.hub == nil {
		return
	}

	msg := Message{
		Type:      MessageType(messageType),
		Data:      data,
		Timestamp: time.Now(),
	}

	n.hub.SendToChannel(channel, msg)
}

// NotifyProgress sends a progress update to a specific user
func (n *Notifier) NotifyProgress(userID uuid.UUID, operationID, step string, progress float64, message string) {
	if n.hub == nil {
		return
	}

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

	n.hub.SendToUser(userID, msg)
}

// NotifyError sends an error message to a specific user
func (n *Notifier) NotifyError(userID uuid.UUID, operationID, message string, err error) {
	if n.hub == nil {
		return
	}

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

	n.hub.SendToUser(userID, msg)
}

// NotifySuccess sends a success message to a specific user
func (n *Notifier) NotifySuccess(userID uuid.UUID, operationID, message string, data interface{}) {
	if n.hub == nil {
		return
	}

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

	n.hub.SendToUser(userID, msg)
}

// NotifyData sends arbitrary data to a specific user
func (n *Notifier) NotifyData(userID uuid.UUID, data interface{}, metadata map[string]interface{}) {
	if n.hub == nil {
		return
	}

	msg := Message{
		Type:      MessageTypeData,
		Data:      data,
		Timestamp: time.Now(),
		Metadata:  metadata,
	}

	n.hub.SendToUser(userID, msg)
}

// NotifyJSON sends a JSON-encoded message to a specific user
func (n *Notifier) NotifyJSON(userID uuid.UUID, messageType string, jsonData string) error {
	if n.hub == nil {
		return nil
	}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return err
	}

	n.Notify(userID, messageType, data)
	return nil
}

// IsConnected checks if a user has any active WebSocket connections
func (n *Notifier) IsConnected(userID uuid.UUID) bool {
	if n.hub == nil {
		return false
	}
	return n.hub.GetUserClientCount(userID) > 0
}

// GetConnectedUserCount returns the total number of connected users
func (n *Notifier) GetConnectedUserCount() int {
	if n.hub == nil {
		return 0
	}
	return n.hub.GetClientCount()
}
