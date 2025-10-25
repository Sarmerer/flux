package entities

import (
	"time"

	"github.com/google/uuid"
)

type Table struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Schema    string    `json:"schema" db:"schema"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type TableCreateRequest struct {
	Name   string                 `json:"name" validate:"required"`
	Schema map[string]interface{} `json:"schema" validate:"required"`
}

type TableUpdateRequest struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
}

type TableResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Schema    string    `json:"schema"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
