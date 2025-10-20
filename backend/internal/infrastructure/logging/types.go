package logging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LogEntry represents a single log entry in the system
type LogEntry struct {
	ID         uuid.UUID              `json:"id"`
	Level      LogLevel               `json:"level"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	ProjectID  *uuid.UUID             `json:"project_id,omitempty"`
	DatabaseID *uuid.UUID             `json:"database_id,omitempty"`
	TableID    *uuid.UUID             `json:"table_id,omitempty"`
	WorkflowID *uuid.UUID             `json:"workflow_id,omitempty"`
	UserID     *uuid.UUID             `json:"user_id,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Operation  string                 `json:"operation,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
}

// LogStorage defines the interface for storing log entries
type LogStorage interface {
	// Store stores a log entry
	Store(ctx context.Context, entry *LogEntry) error

	// Query retrieves log entries based on filters
	Query(ctx context.Context, filter *LogFilter) ([]*LogEntry, error)

	// Count returns the number of log entries matching the filter
	Count(ctx context.Context, filter *LogFilter) (int64, error)

	// Delete deletes log entries older than the specified duration
	DeleteOlderThan(ctx context.Context, duration time.Duration) (int64, error)
}

// LogStreamer defines the interface for streaming log entries to clients
type LogStreamer interface {
	// Stream streams a log entry to connected clients
	Stream(entry *LogEntry) error

	// Subscribe subscribes a client to log streams
	Subscribe(clientID string, filter *LogFilter) error

	// Unsubscribe unsubscribes a client from log streams
	Unsubscribe(clientID string) error
}

// LogFilter defines filters for querying log entries
type LogFilter struct {
	ProjectID  *uuid.UUID `json:"project_id,omitempty"`
	DatabaseID *uuid.UUID `json:"database_id,omitempty"`
	TableID    *uuid.UUID `json:"table_id,omitempty"`
	WorkflowID *uuid.UUID `json:"workflow_id,omitempty"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Level      *LogLevel  `json:"level,omitempty"`
	Operation  *string    `json:"operation,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
	OrderBy    string     `json:"order_by,omitempty"` // timestamp_asc, timestamp_desc
}

// DefaultFilter returns a filter with sensible defaults
func DefaultFilter() *LogFilter {
	return &LogFilter{
		Limit:   100,
		Offset:  0,
		OrderBy: "timestamp_desc",
	}
}
