package entities

import (
	"time"

	"github.com/google/uuid"
)

// Table represents a dynamic table in a project database
type Table struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Schema    string    `json:"schema" db:"schema"` // JSON string of table schema
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableCreateRequest represents the data needed to create a new table
type TableCreateRequest struct {
	Name   string                 `json:"name" validate:"required"`
	Schema map[string]interface{} `json:"schema" validate:"required"`
}

// TableUpdateRequest represents the data needed to update a table
type TableUpdateRequest struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
}

// TableResponse represents the table data returned in API responses
type TableResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Schema    string    `json:"schema"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts a Table entity to TableResponse
func (t *Table) ToResponse() TableResponse {
	return TableResponse{
		ID:        t.ID,
		ProjectID: t.ProjectID,
		Name:      t.Name,
		Schema:    t.Schema,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
