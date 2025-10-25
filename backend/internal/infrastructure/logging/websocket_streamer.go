package logging

import (
	"fmt"
	"sync"

	"github.com/flow/internal/infrastructure/realtime"
)

type WebSocketLogStreamer struct {
	hub           *realtime.Hub
	subscriptions map[string]*LogFilter
	mu            sync.RWMutex
}

func NewWebSocketLogStreamer(hub *realtime.Hub) *WebSocketLogStreamer {
	return &WebSocketLogStreamer{
		hub:           hub,
		subscriptions: make(map[string]*LogFilter),
	}
}

func (s *WebSocketLogStreamer) Stream(entry *LogEntry) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg := realtime.Message{
		Type:      "log",
		Data:      entry,
		Timestamp: entry.Timestamp,
	}

	for clientID, filter := range s.subscriptions {
		if s.matchesFilter(entry, filter) {

			s.hub.SendToClient(clientID, msg)
		}
	}

	s.hub.SendToChannel("logs:all", msg)

	if entry.ProjectID != nil {
		channel := fmt.Sprintf("logs:project:%s", entry.ProjectID.String())
		s.hub.SendToChannel(channel, msg)
	}

	if entry.DatabaseID != nil {
		channel := fmt.Sprintf("logs:database:%s", entry.DatabaseID.String())
		s.hub.SendToChannel(channel, msg)
	}

	if entry.TableID != nil {
		channel := fmt.Sprintf("logs:table:%s", entry.TableID.String())
		s.hub.SendToChannel(channel, msg)
	}

	if entry.WorkflowID != nil {
		channel := fmt.Sprintf("logs:workflow:%s", entry.WorkflowID.String())
		s.hub.SendToChannel(channel, msg)
	}

	return nil
}

func (s *WebSocketLogStreamer) Subscribe(clientID string, filter *LogFilter) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if filter == nil {
		filter = DefaultFilter()
	}

	s.subscriptions[clientID] = filter

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

	if filter.ProjectID == nil && filter.DatabaseID == nil && filter.TableID == nil && filter.WorkflowID == nil {
		s.hub.SubscribeClient(clientID, "logs:all")
	}

	return nil
}

func (s *WebSocketLogStreamer) Unsubscribe(clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if filter, exists := s.subscriptions[clientID]; exists {

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

		s.hub.UnsubscribeClient(clientID, "logs:all")

		delete(s.subscriptions, clientID)
	}

	return nil
}

func (s *WebSocketLogStreamer) matchesFilter(entry *LogEntry, filter *LogFilter) bool {
	if filter == nil {
		return true
	}

	if filter.ProjectID != nil {
		if entry.ProjectID == nil || *entry.ProjectID != *filter.ProjectID {
			return false
		}
	}

	if filter.DatabaseID != nil {
		if entry.DatabaseID == nil || *entry.DatabaseID != *filter.DatabaseID {
			return false
		}
	}

	if filter.TableID != nil {
		if entry.TableID == nil || *entry.TableID != *filter.TableID {
			return false
		}
	}

	if filter.WorkflowID != nil {
		if entry.WorkflowID == nil || *entry.WorkflowID != *filter.WorkflowID {
			return false
		}
	}

	if filter.UserID != nil {
		if entry.UserID == nil || *entry.UserID != *filter.UserID {
			return false
		}
	}

	if filter.Level != nil {
		if entry.Level != *filter.Level {
			return false
		}
	}

	if filter.Operation != nil {
		if entry.Operation != *filter.Operation {
			return false
		}
	}

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
