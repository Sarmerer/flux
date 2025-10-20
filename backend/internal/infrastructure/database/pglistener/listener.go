package pglistener

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lib/pq"
)

// EventType represents the type of database event
type EventType string

const (
	EventTypeInsert EventType = "INSERT"
	EventTypeUpdate EventType = "UPDATE"
	EventTypeDelete EventType = "DELETE"
)

// Event represents a database change event
type Event struct {
	Type      EventType              `json:"type"`
	Table     string                 `json:"table"`
	Schema    string                 `json:"schema"`
	Timestamp time.Time              `json:"timestamp"`
	OldData   map[string]interface{} `json:"old_data,omitempty"`
	NewData   map[string]interface{} `json:"new_data,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// EventHandler is a function that processes database events
type EventHandler func(event Event)

// Listener manages PostgreSQL LISTEN/NOTIFY subscriptions
type Listener struct {
	db              *sql.DB
	listener        *pq.Listener
	handlers        map[string][]EventHandler // channel -> handlers
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	reconnectDelay  time.Duration
	maxReconnectDelay time.Duration
}

// Config holds configuration for the listener
type Config struct {
	DatabaseURL       string
	MinReconnectDelay time.Duration
	MaxReconnectDelay time.Duration
}

// NewListener creates a new PostgreSQL listener
func NewListener(db *sql.DB, config Config) (*Listener, error) {
	ctx, cancel := context.WithCancel(context.Background())

	minDelay := config.MinReconnectDelay
	if minDelay == 0 {
		minDelay = 10 * time.Second
	}

	maxDelay := config.MaxReconnectDelay
	if maxDelay == 0 {
		maxDelay = time.Minute
	}

	eventCallback := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[pglistener] Listener event error: %v", err)
		}
		switch ev {
		case pq.ListenerEventConnected:
			log.Println("[pglistener] Connected to PostgreSQL notification channel")
		case pq.ListenerEventDisconnected:
			log.Println("[pglistener] Disconnected from PostgreSQL notification channel")
		case pq.ListenerEventReconnected:
			log.Println("[pglistener] Reconnected to PostgreSQL notification channel")
		case pq.ListenerEventConnectionAttemptFailed:
			log.Printf("[pglistener] Connection attempt failed: %v", err)
		}
	}

	pqListener := pq.NewListener(
		config.DatabaseURL,
		minDelay,
		maxDelay,
		eventCallback,
	)

	l := &Listener{
		db:                db,
		listener:          pqListener,
		handlers:          make(map[string][]EventHandler),
		ctx:               ctx,
		cancel:            cancel,
		reconnectDelay:    minDelay,
		maxReconnectDelay: maxDelay,
	}

	return l, nil
}

// Subscribe adds a handler for a specific channel
func (l *Listener) Subscribe(channel string, handler EventHandler) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Add handler to the list
	l.handlers[channel] = append(l.handlers[channel], handler)

	// Listen on the channel if this is the first handler
	if len(l.handlers[channel]) == 1 {
		err := l.listener.Listen(channel)
		if err != nil {
			return fmt.Errorf("failed to listen on channel %s: %w", channel, err)
		}
		log.Printf("[pglistener] Subscribed to channel: %s", channel)
	}

	return nil
}

// Unsubscribe removes a specific handler from a channel
// Note: This is a simplified version. For production, you'd want to track handler identity
func (l *Listener) Unsubscribe(channel string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.handlers, channel)

	err := l.listener.Unlisten(channel)
	if err != nil {
		return fmt.Errorf("failed to unlisten from channel %s: %w", channel, err)
	}

	log.Printf("[pglistener] Unsubscribed from channel: %s", channel)
	return nil
}

// Start begins listening for notifications
func (l *Listener) Start() error {
	log.Println("[pglistener] Starting listener...")

	go l.listen()

	return nil
}

// listen is the main event loop
func (l *Listener) listen() {
	for {
		select {
		case <-l.ctx.Done():
			log.Println("[pglistener] Listener context cancelled, shutting down...")
			return

		case notification := <-l.listener.Notify:
			if notification == nil {
				// Nil notification can occur during reconnection
				continue
			}

			l.handleNotification(notification)

		case <-time.After(90 * time.Second):
			// Send a ping to check connection health
			go func() {
				err := l.listener.Ping()
				if err != nil {
					log.Printf("[pglistener] Ping failed: %v", err)
				}
			}()
		}
	}
}

// handleNotification processes an incoming notification
func (l *Listener) handleNotification(n *pq.Notification) {
	l.mu.RLock()
	handlers, exists := l.handlers[n.Channel]
	l.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		log.Printf("[pglistener] No handlers for channel: %s", n.Channel)
		return
	}

	// Parse the notification payload
	var event Event
	err := json.Unmarshal([]byte(n.Extra), &event)
	if err != nil {
		log.Printf("[pglistener] Failed to parse notification payload: %v", err)
		return
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Call all handlers for this channel
	for _, handler := range handlers {
		go handler(event)
	}
}

// Stop gracefully shuts down the listener
func (l *Listener) Stop() error {
	log.Println("[pglistener] Stopping listener...")

	l.cancel()

	err := l.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	log.Println("[pglistener] Listener stopped")
	return nil
}

// Notify sends a notification to a channel
// This is a helper method for sending notifications from application code
func (l *Listener) Notify(channel string, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	query := fmt.Sprintf("NOTIFY %s, '%s'", channel, string(payload))
	_, err = l.db.ExecContext(l.ctx, query)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

// GetActiveChannels returns a list of channels currently being listened to
func (l *Listener) GetActiveChannels() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	channels := make([]string, 0, len(l.handlers))
	for channel := range l.handlers {
		channels = append(channels, channel)
	}
	return channels
}
