package domain

import "time"

type ColumnType string

const (
	ColumnText      ColumnType = "text"
	ColumnInteger   ColumnType = "integer"
	ColumnBoolean   ColumnType = "boolean"
	ColumnTimestamp ColumnType = "timestamp"
)

type Column struct {
	ID           int64      `json:"id"`
	TableID      int64      `json:"table_id"`
	Name         string     `json:"name"`
	Type         ColumnType `json:"type"`
	Required     bool       `json:"required"`
	Unique       bool       `json:"unique"`
	DefaultValue *string    `json:"default_value"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
