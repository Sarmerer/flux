package postgres

import (
	"context"
	"fmt"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/domain/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) repositories.ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *entities.Project) error {
	return r.CreateTx(ctx, r.db, project)
}

func (r *ProjectRepository) CreateTx(ctx context.Context, tx types.Executor, project *entities.Project) error {
	query := `
		INSERT INTO projects (id, name, description, owner_id, database_url, api_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := tx.Exec(ctx, query, project.ID, project.Name, project.Description, project.OwnerID, project.DatabaseURL, project.APIKey, project.CreatedAt, project.UpdatedAt)
	return err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	query := `
		SELECT id, name, description, owner_id, database_url, api_key, created_at, updated_at
		FROM projects WHERE id = $1
	`
	var project entities.Project
	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.DatabaseURL, &project.APIKey, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entities.Project, error) {
	query := `
		SELECT id, name, description, owner_id, database_url, api_key, created_at, updated_at
		FROM projects WHERE owner_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*entities.Project
	for rows.Next() {
		var project entities.Project
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.DatabaseURL, &project.APIKey, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *entities.Project) error {
	return r.UpdateTx(ctx, r.db, project)
}

func (r *ProjectRepository) UpdateTx(ctx context.Context, tx types.Executor, project *entities.Project) error {
	query := `
		UPDATE projects
		SET name = $2, description = $3, database_url = $4, api_key = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := tx.Exec(ctx, query, project.ID, project.Name, project.Description, project.DatabaseURL, project.APIKey, project.UpdatedAt)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteTx(ctx, r.db, id)
}

func (r *ProjectRepository) DeleteTx(ctx context.Context, tx types.Executor, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	return err
}

func (r *ProjectRepository) List(ctx context.Context, limit, offset int) ([]*entities.Project, error) {
	query := `
		SELECT id, name, description, owner_id, database_url, api_key, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*entities.Project
	for rows.Next() {
		var project entities.Project
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.DatabaseURL, &project.APIKey, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	return projects, nil
}
