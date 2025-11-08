package entities

import (
	"time"

	"github.com/google/uuid"
)

type Table struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TableCreateRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

type TableUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TableResponse struct {
	ID          uuid.UUID          `json:"id"`
	ProjectID   uuid.UUID          `json:"project_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Schema      *TableSchemaInfo   `json:"schema,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func (t *Table) ToResponse() TableResponse {
	return TableResponse{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Name:        t.Name,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
