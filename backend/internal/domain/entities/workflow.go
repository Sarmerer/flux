package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Workflow represents an automation workflow
type Workflow struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	ProjectID   uuid.UUID       `json:"project_id" db:"project_id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	Trigger     WorkflowTrigger `json:"trigger" db:"trigger"`
	Actions     WorkflowActions `json:"actions" db:"actions"`
	IsActive    bool            `json:"is_active" db:"is_active"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// WorkflowTrigger represents a workflow trigger configuration
type WorkflowTrigger struct {
	Type       string                 `json:"type"` // on_row_created, on_row_updated, on_row_deleted, scheduled, webhook
	TableName  string                 `json:"table_name,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Schedule   string                 `json:"schedule,omitempty"` // Cron expression
}

// WorkflowAction represents a single action in a workflow
type WorkflowAction struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"` // send_webhook, send_email, update_row, create_row, delete_row
	Config map[string]interface{} `json:"config"`
	Order  int                    `json:"order"`
}

// WorkflowActions is a slice of WorkflowAction that implements sql.Scanner and driver.Valuer
type WorkflowActions []WorkflowAction

// Scan implements sql.Scanner for WorkflowActions
func (wa *WorkflowActions) Scan(value interface{}) error {
	if value == nil {
		*wa = WorkflowActions{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, wa)
}

// Value implements driver.Valuer for WorkflowActions
func (wa WorkflowActions) Value() (driver.Value, error) {
	if wa == nil {
		return nil, nil
	}
	return json.Marshal(wa)
}

// Scan implements sql.Scanner for WorkflowTrigger
func (wt *WorkflowTrigger) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, wt)
}

// Value implements driver.Valuer for WorkflowTrigger
func (wt WorkflowTrigger) Value() (driver.Value, error) {
	return json.Marshal(wt)
}

// WorkflowCreateRequest represents the data needed to create a workflow
type WorkflowCreateRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Trigger     WorkflowTrigger `json:"trigger" validate:"required"`
	Actions     []WorkflowAction `json:"actions" validate:"required,min=1"`
}

// WorkflowUpdateRequest represents the data needed to update a workflow
type WorkflowUpdateRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Trigger     *WorkflowTrigger `json:"trigger"`
	Actions     []WorkflowAction `json:"actions"`
	IsActive    *bool            `json:"is_active"`
}

// WorkflowResponse represents the workflow data returned in API responses
type WorkflowResponse struct {
	ID          uuid.UUID        `json:"id"`
	ProjectID   uuid.UUID        `json:"project_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Trigger     WorkflowTrigger  `json:"trigger"`
	Actions     []WorkflowAction `json:"actions"`
	IsActive    bool             `json:"is_active"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ToResponse converts a Workflow entity to WorkflowResponse
func (w *Workflow) ToResponse() WorkflowResponse {
	return WorkflowResponse{
		ID:          w.ID,
		ProjectID:   w.ProjectID,
		Name:        w.Name,
		Description: w.Description,
		Trigger:     w.Trigger,
		Actions:     w.Actions,
		IsActive:    w.IsActive,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}
