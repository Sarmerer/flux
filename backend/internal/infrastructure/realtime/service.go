package realtime

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/flow/internal/infrastructure/database/pglistener"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Service integrates pglistener and realtime hub to provide full realtime functionality
type Service struct {
	hub      *Hub
	notifier *Notifier
	listener *pglistener.Listener
	triggerMgr *pglistener.TriggerManager
}

// Config holds configuration for the realtime service
type Config struct {
	DatabaseURL       string
	JWTSecret         string
}

// NewService creates a new realtime service
func NewService(databaseURL string, config Config) (*Service, error) {
	// Create hub
	hub := NewHub()
	notifier := NewNotifier(hub)

	// Open SQL connection for pglistener
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create listener
	listenerConfig := pglistener.Config{
		DatabaseURL: config.DatabaseURL,
	}

	listener, err := pglistener.NewListener(db, listenerConfig)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create pglistener: %w", err)
	}

	// Create trigger manager
	triggerMgr := pglistener.NewTriggerManager(db)

	return &Service{
		hub:        hub,
		notifier:   notifier,
		listener:   listener,
		triggerMgr: triggerMgr,
	}, nil
}

// Start starts the realtime service
func (s *Service) Start(ctx context.Context) error {
	log.Println("[realtime] Starting realtime service...")

	// Start the hub
	go s.hub.Run(ctx)

	// Start the listener
	err := s.listener.Start()
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	log.Println("[realtime] Realtime service started")
	return nil
}

// Stop stops the realtime service
func (s *Service) Stop() error {
	log.Println("[realtime] Stopping realtime service...")

	err := s.listener.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop listener: %w", err)
	}

	log.Println("[realtime] Realtime service stopped")
	return nil
}

// GetHub returns the WebSocket hub
func (s *Service) GetHub() *Hub {
	return s.hub
}

// GetNotifier returns the notifier
func (s *Service) GetNotifier() *Notifier {
	return s.notifier
}

// GetListener returns the pglistener
func (s *Service) GetListener() *pglistener.Listener {
	return s.listener
}

// GetTriggerManager returns the trigger manager
func (s *Service) GetTriggerManager() *pglistener.TriggerManager {
	return s.triggerMgr
}

// EnableRealtimeForTable enables realtime updates for a specific table
func (s *Service) EnableRealtimeForTable(ctx context.Context, schema, table string) error {
	channel := pglistener.CreateRealtimeChannel(schema, table)

	// Create database trigger
	err := s.triggerMgr.EnableTriggersForTable(ctx, schema, table, channel)
	if err != nil {
		return fmt.Errorf("failed to enable triggers: %w", err)
	}

	// Subscribe to the channel and forward events to WebSocket
	handler := func(event pglistener.Event) {
		// Broadcast the event to all subscribed clients
		msg := Message{
			Type:      MessageTypeRealtimeData,
			Timestamp: event.Timestamp,
			Data: map[string]interface{}{
				"event_type": string(event.Type),
				"table":      event.Table,
				"schema":     event.Schema,
				"old_data":   event.OldData,
				"new_data":   event.NewData,
			},
			Metadata: event.Metadata,
		}

		// Send to channel subscribers
		s.hub.SendToChannel(channel, msg)
	}

	err = s.listener.Subscribe(channel, handler)
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel: %w", err)
	}

	log.Printf("[realtime] Enabled realtime for %s.%s on channel %s", schema, table, channel)
	return nil
}

// DisableRealtimeForTable disables realtime updates for a specific table
func (s *Service) DisableRealtimeForTable(ctx context.Context, schema, table string) error {
	channel := pglistener.CreateRealtimeChannel(schema, table)

	// Remove database trigger
	err := s.triggerMgr.DisableTriggersForTable(ctx, schema, table)
	if err != nil {
		return fmt.Errorf("failed to disable triggers: %w", err)
	}

	// Unsubscribe from the channel
	err = s.listener.Unsubscribe(channel)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from channel: %w", err)
	}

	log.Printf("[realtime] Disabled realtime for %s.%s", schema, table)
	return nil
}

// SubscribeClientToTable subscribes a WebSocket client to table changes
func (s *Service) SubscribeClientToTable(clientID string, schema, table string) {
	channel := pglistener.CreateRealtimeChannel(schema, table)
	s.hub.SubscribeClient(clientID, channel)
}

// UnsubscribeClientFromTable unsubscribes a WebSocket client from table changes
func (s *Service) UnsubscribeClientFromTable(clientID string, schema, table string) {
	channel := pglistener.CreateRealtimeChannel(schema, table)
	s.hub.UnsubscribeClient(clientID, channel)
}

// NotifyDatabaseChange sends a database change notification
// This is useful for manual notifications outside the trigger system
func (s *Service) NotifyDatabaseChange(ctx context.Context, schema, table string, eventType pglistener.EventType, oldData, newData map[string]interface{}) error {
	channel := pglistener.CreateRealtimeChannel(schema, table)

	event := pglistener.Event{
		Type:    eventType,
		Table:   table,
		Schema:  schema,
		OldData: oldData,
		NewData: newData,
	}

	return s.listener.Notify(channel, event)
}

// SendProgressToUser sends progress updates to a specific user
func (s *Service) SendProgressToUser(userID uuid.UUID, operationID, step string, progress float64, message string) {
	s.notifier.NotifyProgress(userID, operationID, step, progress, message)
}

// SendErrorToUser sends error messages to a specific user
func (s *Service) SendErrorToUser(userID uuid.UUID, operationID, message string, err error) {
	s.notifier.NotifyError(userID, operationID, message, err)
}

// SendSuccessToUser sends success messages to a specific user
func (s *Service) SendSuccessToUser(userID uuid.UUID, operationID, message string, data interface{}) {
	s.notifier.NotifySuccess(userID, operationID, message, data)
}

// GetStats returns realtime service statistics
func (s *Service) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"connected_clients":   s.hub.GetClientCount(),
		"active_channels":     s.listener.GetActiveChannels(),
		"channel_count":       len(s.listener.GetActiveChannels()),
	}
}
