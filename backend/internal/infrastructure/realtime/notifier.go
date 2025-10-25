package realtime

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Notifier struct {
	hub *Hub
}

func NewNotifier(hub *Hub) *Notifier {
	return &Notifier{hub: hub}
}

func (n *Notifier) Notify(userID uuid.UUID, messageType string, data interface{}) {
	if n.hub == nil {
		return
	}

	msg := Message{
		Type:      MessageType(messageType),
		Data:      data,
		Timestamp: time.Now(),
	}

	n.hub.SendToUser(userID, msg)
}

func (n *Notifier) NotifyWithID(userID uuid.UUID, messageType string, operationID string, data interface{}) {
	if n.hub == nil {
		return
	}

	msg := Message{
		Type:      MessageType(messageType),
		ID:        operationID,
		Data:      data,
		Timestamp: time.Now(),
	}

	n.hub.SendToUser(userID, msg)
}

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

func (n *Notifier) IsConnected(userID uuid.UUID) bool {
	if n.hub == nil {
		return false
	}
	return n.hub.GetUserClientCount(userID) > 0
}

func (n *Notifier) GetConnectedUserCount() int {
	if n.hub == nil {
		return 0
	}
	return n.hub.GetClientCount()
}
