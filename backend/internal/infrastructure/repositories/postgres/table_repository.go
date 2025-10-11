package postgres

import (
	"context"
	"fmt"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TableRepository implements the TableRepository interface using PostgreSQL
type TableRepository struct {
	db *pgxpool.Pool
}

// NewTableRepository creates a new TableRepository
func NewTableRepository(db *pgxpool.Pool) repositories.TableRepository {
	return &TableRepository{db: db}
}

// Create creates a new table
func (r *TableRepository) Create(ctx context.Context, table *entities.Table) error {
	query := `
		INSERT INTO tables (id, project_id, name, schema, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, table.ID, table.ProjectID, table.Name, table.Schema, table.CreatedAt, table.UpdatedAt)
	return err
}

// GetByID retrieves a table by ID
func (r *TableRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Table, error) {
	query := `
		SELECT id, project_id, name, schema, created_at, updated_at
		FROM tables WHERE id = $1
	`
	var table entities.Table
	err := r.db.QueryRow(ctx, query, id).Scan(
		&table.ID, &table.ProjectID, &table.Name, &table.Schema, &table.CreatedAt, &table.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("table not found")
		}
		return nil, err
	}
	return &table, nil
}

// GetByProjectID retrieves tables by project ID
func (r *TableRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Table, error) {
	query := `
		SELECT id, project_id, name, schema, created_at, updated_at
		FROM tables WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*entities.Table
	for rows.Next() {
		var table entities.Table
		err := rows.Scan(&table.ID, &table.ProjectID, &table.Name, &table.Schema, &table.CreatedAt, &table.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tables = append(tables, &table)
	}

	return tables, nil
}

// Update updates a table
func (r *TableRepository) Update(ctx context.Context, table *entities.Table) error {
	query := `
		UPDATE tables 
		SET name = $2, schema = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, table.ID, table.Name, table.Schema, table.UpdatedAt)
	return err
}

// Delete deletes a table
func (r *TableRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tables WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// List retrieves tables with pagination
func (r *TableRepository) List(ctx context.Context, limit, offset int) ([]*entities.Table, error) {
	query := `
		SELECT id, project_id, name, schema, created_at, updated_at
		FROM tables 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*entities.Table
	for rows.Next() {
		var table entities.Table
		err := rows.Scan(&table.ID, &table.ProjectID, &table.Name, &table.Schema, &table.CreatedAt, &table.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tables = append(tables, &table)
	}

	return tables, nil
}
