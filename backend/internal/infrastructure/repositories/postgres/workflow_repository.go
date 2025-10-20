package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WorkflowRepository implements the WorkflowRepository interface using PostgreSQL
type WorkflowRepository struct {
	db *pgxpool.Pool
}

// NewWorkflowRepository creates a new WorkflowRepository
func NewWorkflowRepository(db *pgxpool.Pool) repositories.WorkflowRepository {
	return &WorkflowRepository{db: db}
}

// Create creates a new workflow
func (r *WorkflowRepository) Create(ctx context.Context, workflow *entities.Workflow) error {
	query := `
		INSERT INTO workflows (id, project_id, name, description, trigger, actions, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		workflow.ID,
		workflow.ProjectID,
		workflow.Name,
		workflow.Description,
		workflow.Trigger,
		workflow.Actions,
		workflow.IsActive,
		workflow.CreatedAt,
		workflow.UpdatedAt,
	)
	return err
}

// GetByID retrieves a workflow by ID
func (r *WorkflowRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Workflow, error) {
	query := `
		SELECT id, project_id, name, description, trigger, actions, is_active, created_at, updated_at
		FROM workflows WHERE id = $1
	`
	var workflow entities.Workflow
	err := r.db.QueryRow(ctx, query, id).Scan(
		&workflow.ID,
		&workflow.ProjectID,
		&workflow.Name,
		&workflow.Description,
		&workflow.Trigger,
		&workflow.Actions,
		&workflow.IsActive,
		&workflow.CreatedAt,
		&workflow.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("workflow not found")
		}
		return nil, err
	}
	return &workflow, nil
}

// GetByProjectID retrieves workflows by project ID
func (r *WorkflowRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Workflow, error) {
	query := `
		SELECT id, project_id, name, description, trigger, actions, is_active, created_at, updated_at
		FROM workflows WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []*entities.Workflow
	for rows.Next() {
		var workflow entities.Workflow
		err := rows.Scan(
			&workflow.ID,
			&workflow.ProjectID,
			&workflow.Name,
			&workflow.Description,
			&workflow.Trigger,
			&workflow.Actions,
			&workflow.IsActive,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, &workflow)
	}

	return workflows, nil
}

// Update updates a workflow
func (r *WorkflowRepository) Update(ctx context.Context, workflow *entities.Workflow) error {
	query := `
		UPDATE workflows
		SET name = $2, description = $3, trigger = $4, actions = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query,
		workflow.ID,
		workflow.Name,
		workflow.Description,
		workflow.Trigger,
		workflow.Actions,
		workflow.IsActive,
		workflow.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("workflow not found")
	}

	return nil
}

// Delete deletes a workflow
func (r *WorkflowRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM workflows WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("workflow not found")
	}

	return nil
}

// List retrieves workflows with pagination
func (r *WorkflowRepository) List(ctx context.Context, limit, offset int) ([]*entities.Workflow, error) {
	query := `
		SELECT id, project_id, name, description, trigger, actions, is_active, created_at, updated_at
		FROM workflows
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []*entities.Workflow
	for rows.Next() {
		var workflow entities.Workflow
		err := rows.Scan(
			&workflow.ID,
			&workflow.ProjectID,
			&workflow.Name,
			&workflow.Description,
			&workflow.Trigger,
			&workflow.Actions,
			&workflow.IsActive,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, &workflow)
	}

	return workflows, nil
}

// ToggleActive toggles the active state of a workflow
func (r *WorkflowRepository) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	query := `
		UPDATE workflows
		SET is_active = $2, updated_at = $3
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, id, isActive, time.Now())
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("workflow not found")
	}

	return nil
}
