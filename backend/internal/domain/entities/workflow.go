package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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

type WorkflowTrigger struct {
	Type       string                 `json:"type"`
	TableName  string                 `json:"table_name,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Schedule   string                 `json:"schedule,omitempty"`
}

type WorkflowAction struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
	Order  int                    `json:"order"`
}

type WorkflowActions []WorkflowAction

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

func (wa WorkflowActions) Value() (driver.Value, error) {
	if wa == nil {
		return nil, nil
	}
	return json.Marshal(wa)
}

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

func (wt WorkflowTrigger) Value() (driver.Value, error) {
	return json.Marshal(wt)
}

type WorkflowCreateRequest struct {
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description"`
	Trigger     WorkflowTrigger  `json:"trigger" validate:"required"`
	Actions     []WorkflowAction `json:"actions" validate:"required,min=1"`
}

type WorkflowUpdateRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Trigger     *WorkflowTrigger `json:"trigger"`
	Actions     []WorkflowAction `json:"actions"`
	IsActive    *bool            `json:"is_active"`
}

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

func (w *Workflow) ToResponse() WorkflowResponse {
	actions := w.Actions
	if actions == nil {
		actions = []WorkflowAction{}
	}
	return WorkflowResponse{
		ID:          w.ID,
		ProjectID:   w.ProjectID,
		Name:        w.Name,
		Description: w.Description,
		Trigger:     w.Trigger,
		Actions:     actions,
		IsActive:    w.IsActive,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}
