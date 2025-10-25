package logging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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

type LogStorage interface {
	Store(ctx context.Context, entry *LogEntry) error

	Query(ctx context.Context, filter *LogFilter) ([]*LogEntry, error)

	Count(ctx context.Context, filter *LogFilter) (int64, error)

	DeleteOlderThan(ctx context.Context, duration time.Duration) (int64, error)
}

type LogStreamer interface {
	Stream(entry *LogEntry) error

	Subscribe(clientID string, filter *LogFilter) error

	Unsubscribe(clientID string) error
}

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
	OrderBy    string     `json:"order_by,omitempty"`
}

func DefaultFilter() *LogFilter {
	return &LogFilter{
		Limit:   100,
		Offset:  0,
		OrderBy: "timestamp_desc",
	}
}
