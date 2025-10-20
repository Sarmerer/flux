package logging

import (
	"fmt"
	"sync"

	"github.com/flow/internal/infrastructure/realtime"
)

// WebSocketLogStreamer implements LogStreamer using WebSocket
type WebSocketLogStreamer struct {
	hub           *realtime.Hub
	subscriptions map[string]*LogFilter // clientID -> filter
	mu            sync.RWMutex
}

// NewWebSocketLogStreamer creates a new WebSocket log streamer
func NewWebSocketLogStreamer(hub *realtime.Hub) *WebSocketLogStreamer {
	return &WebSocketLogStreamer{
		hub:           hub,
		subscriptions: make(map[string]*LogFilter),
	}
}

// Stream streams a log entry to subscribed WebSocket clients
func (s *WebSocketLogStreamer) Stream(entry *LogEntry) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create WebSocket message
	msg := realtime.Message{
		Type:      "log",
		Data:      entry,
		Timestamp: entry.Timestamp,
	}

	// Check which clients should receive this log
	for clientID, filter := range s.subscriptions {
		if s.matchesFilter(entry, filter) {
			// Send to specific client
			s.hub.SendToClient(clientID, msg)
		}
	}

	// Also send to global log channel (for admin views)
	s.hub.SendToChannel("logs:all", msg)

	// Send to project-specific channel if project_id exists
	if entry.ProjectID != nil {
		channel := fmt.Sprintf("logs:project:%s", entry.ProjectID.String())
		s.hub.SendToChannel(channel, msg)
	}

	// Send to database-specific channel if database_id exists
	if entry.DatabaseID != nil {
		channel := fmt.Sprintf("logs:database:%s", entry.DatabaseID.String())
		s.hub.SendToChannel(channel, msg)
	}

	// Send to table-specific channel if table_id exists
	if entry.TableID != nil {
		channel := fmt.Sprintf("logs:table:%s", entry.TableID.String())
		s.hub.SendToChannel(channel, msg)
	}

	// Send to workflow-specific channel if workflow_id exists
	if entry.WorkflowID != nil {
		channel := fmt.Sprintf("logs:workflow:%s", entry.WorkflowID.String())
		s.hub.SendToChannel(channel, msg)
	}

	return nil
}

// Subscribe subscribes a client to log streams with a filter
func (s *WebSocketLogStreamer) Subscribe(clientID string, filter *LogFilter) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if filter == nil {
		filter = DefaultFilter()
	}

	s.subscriptions[clientID] = filter

	// Subscribe to relevant channels
	if filter.ProjectID != nil {
		channel := fmt.Sprintf("logs:project:%s", filter.ProjectID.String())
		s.hub.SubscribeClient(clientID, channel)
	}

	if filter.DatabaseID != nil {
		channel := fmt.Sprintf("logs:database:%s", filter.DatabaseID.String())
		s.hub.SubscribeClient(clientID, channel)
	}

	if filter.TableID != nil {
		channel := fmt.Sprintf("logs:table:%s", filter.TableID.String())
		s.hub.SubscribeClient(clientID, channel)
	}

	if filter.WorkflowID != nil {
		channel := fmt.Sprintf("logs:workflow:%s", filter.WorkflowID.String())
		s.hub.SubscribeClient(clientID, channel)
	}

	// Subscribe to global logs if no specific filter
	if filter.ProjectID == nil && filter.DatabaseID == nil && filter.TableID == nil && filter.WorkflowID == nil {
		s.hub.SubscribeClient(clientID, "logs:all")
	}

	return nil
}

// Unsubscribe unsubscribes a client from log streams
func (s *WebSocketLogStreamer) Unsubscribe(clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove subscription
	if filter, exists := s.subscriptions[clientID]; exists {
		// Unsubscribe from channels
		if filter.ProjectID != nil {
			channel := fmt.Sprintf("logs:project:%s", filter.ProjectID.String())
			s.hub.UnsubscribeClient(clientID, channel)
		}

		if filter.DatabaseID != nil {
			channel := fmt.Sprintf("logs:database:%s", filter.DatabaseID.String())
			s.hub.UnsubscribeClient(clientID, channel)
		}

		if filter.TableID != nil {
			channel := fmt.Sprintf("logs:table:%s", filter.TableID.String())
			s.hub.UnsubscribeClient(clientID, channel)
		}

		if filter.WorkflowID != nil {
			channel := fmt.Sprintf("logs:workflow:%s", filter.WorkflowID.String())
			s.hub.UnsubscribeClient(clientID, channel)
		}

		// Unsubscribe from global logs
		s.hub.UnsubscribeClient(clientID, "logs:all")

		delete(s.subscriptions, clientID)
	}

	return nil
}

// matchesFilter checks if a log entry matches a filter
func (s *WebSocketLogStreamer) matchesFilter(entry *LogEntry, filter *LogFilter) bool {
	if filter == nil {
		return true
	}

	// Check project ID
	if filter.ProjectID != nil {
		if entry.ProjectID == nil || *entry.ProjectID != *filter.ProjectID {
			return false
		}
	}

	// Check database ID
	if filter.DatabaseID != nil {
		if entry.DatabaseID == nil || *entry.DatabaseID != *filter.DatabaseID {
			return false
		}
	}

	// Check table ID
	if filter.TableID != nil {
		if entry.TableID == nil || *entry.TableID != *filter.TableID {
			return false
		}
	}

	// Check workflow ID
	if filter.WorkflowID != nil {
		if entry.WorkflowID == nil || *entry.WorkflowID != *filter.WorkflowID {
			return false
		}
	}

	// Check user ID
	if filter.UserID != nil {
		if entry.UserID == nil || *entry.UserID != *filter.UserID {
			return false
		}
	}

	// Check level
	if filter.Level != nil {
		if entry.Level != *filter.Level {
			return false
		}
	}

	// Check operation
	if filter.Operation != nil {
		if entry.Operation != *filter.Operation {
			return false
		}
	}

	// Check time range
	if filter.StartTime != nil {
		if entry.Timestamp.Before(*filter.StartTime) {
			return false
		}
	}

	if filter.EndTime != nil {
		if entry.Timestamp.After(*filter.EndTime) {
			return false
		}
	}

	return true
}
