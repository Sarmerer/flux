package realtime

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/flow/internal/infrastructure/database/pglistener"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Service struct {
	hub        *Hub
	notifier   *Notifier
	listener   *pglistener.Listener
	triggerMgr *pglistener.TriggerManager
}

type Config struct {
	DatabaseURL string
	JWTSecret   string
}

func NewService(databaseURL string, config Config) (*Service, error) {

	hub := NewHub()
	notifier := NewNotifier(hub)

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	listenerConfig := pglistener.Config{
		DatabaseURL: config.DatabaseURL,
	}

	listener, err := pglistener.NewListener(db, listenerConfig)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to create pglistener: %w", err)
	}

	triggerMgr := pglistener.NewTriggerManager(db)

	return &Service{
		hub:        hub,
		notifier:   notifier,
		listener:   listener,
		triggerMgr: triggerMgr,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	log.Println("[realtime] Starting realtime service...")

	go s.hub.Run(ctx)

	err := s.listener.Start()
	if err != nil {
		return fmt.Errorf("Failed to start listener: %w", err)
	}

	log.Println("[realtime] Realtime service started")
	return nil
}

func (s *Service) Stop() error {
	log.Println("[realtime] Stopping realtime service...")

	err := s.listener.Stop()
	if err != nil {
		return fmt.Errorf("Failed to stop listener: %w", err)
	}

	log.Println("[realtime] Realtime service stopped")
	return nil
}

func (s *Service) GetHub() *Hub {
	return s.hub
}

func (s *Service) GetNotifier() *Notifier {
	return s.notifier
}

func (s *Service) GetListener() *pglistener.Listener {
	return s.listener
}

func (s *Service) GetTriggerManager() *pglistener.TriggerManager {
	return s.triggerMgr
}

func (s *Service) EnableRealtimeForTable(ctx context.Context, schema, table string) error {
	channel := pglistener.CreateRealtimeChannel(schema, table)

	err := s.triggerMgr.EnableTriggersForTable(ctx, schema, table, channel)
	if err != nil {
		return fmt.Errorf("Failed to enable triggers: %w", err)
	}

	handler := func(event pglistener.Event) {

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

		s.hub.SendToChannel(channel, msg)
	}

	err = s.listener.Subscribe(channel, handler)
	if err != nil {
		return fmt.Errorf("Failed to subscribe to channel: %w", err)
	}

	log.Printf("[realtime] Enabled realtime for %s.%s on channel %s", schema, table, channel)
	return nil
}

func (s *Service) DisableRealtimeForTable(ctx context.Context, schema, table string) error {
	channel := pglistener.CreateRealtimeChannel(schema, table)

	err := s.triggerMgr.DisableTriggersForTable(ctx, schema, table)
	if err != nil {
		return fmt.Errorf("Failed to disable triggers: %w", err)
	}

	err = s.listener.Unsubscribe(channel)
	if err != nil {
		return fmt.Errorf("Failed to unsubscribe from channel: %w", err)
	}

	log.Printf("[realtime] Disabled realtime for %s.%s", schema, table)
	return nil
}

func (s *Service) SubscribeClientToTable(clientID string, schema, table string) {
	channel := pglistener.CreateRealtimeChannel(schema, table)
	s.hub.SubscribeClient(clientID, channel)
}

func (s *Service) UnsubscribeClientFromTable(clientID string, schema, table string) {
	channel := pglistener.CreateRealtimeChannel(schema, table)
	s.hub.UnsubscribeClient(clientID, channel)
}

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

func (s *Service) SendProgressToUser(userID uuid.UUID, operationID, step string, progress float64, message string) {
	s.notifier.NotifyProgress(userID, operationID, step, progress, message)
}

func (s *Service) SendErrorToUser(userID uuid.UUID, operationID, message string, err error) {
	s.notifier.NotifyError(userID, operationID, message, err)
}

func (s *Service) SendSuccessToUser(userID uuid.UUID, operationID, message string, data interface{}) {
	s.notifier.NotifySuccess(userID, operationID, message, data)
}

func (s *Service) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"connected_clients": s.hub.GetClientCount(),
		"active_channels":   s.listener.GetActiveChannels(),
		"channel_count":     len(s.listener.GetActiveChannels()),
	}
}
