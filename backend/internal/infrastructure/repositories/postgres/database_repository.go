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

type DatabaseRepository struct {
	db *pgxpool.Pool
}

func NewDatabaseRepository(db *pgxpool.Pool) repositories.DatabaseRepository {
	return &DatabaseRepository{db: db}
}

func (r *DatabaseRepository) Create(ctx context.Context, database *entities.Database) error {
	query := `
		INSERT INTO databases (id, project_id, name, host, port, username, password, database, ssl_mode, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query, database.ID, database.ProjectID, database.Name, database.Host, database.Port, database.Username, database.Password, database.Database, database.SSLMode, database.CreatedAt, database.UpdatedAt)
	return err
}

func (r *DatabaseRepository) CreateTx(ctx context.Context, tx types.Executor, database *entities.Database) error {
	query := `
		INSERT INTO databases (id, project_id, name, host, port, username, password, database, ssl_mode, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := tx.Exec(ctx, query, database.ID, database.ProjectID, database.Name, database.Host, database.Port, database.Username, database.Password, database.Database, database.SSLMode, database.CreatedAt, database.UpdatedAt)
	return err
}

func (r *DatabaseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Database, error) {
	query := `
		SELECT id, project_id, name, host, port, username, password, database, ssl_mode, created_at, updated_at
		FROM databases WHERE id = $1
	`
	var database entities.Database
	err := r.db.QueryRow(ctx, query, id).Scan(
		&database.ID, &database.ProjectID, &database.Name, &database.Host, &database.Port, &database.Username, &database.Password, &database.Database, &database.SSLMode, &database.CreatedAt, &database.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("database not found")
		}
		return nil, err
	}
	return &database, nil
}

func (r *DatabaseRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.Database, error) {
	query := `
		SELECT id, project_id, name, host, port, username, password, database, ssl_mode, created_at, updated_at
		FROM databases WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []*entities.Database
	for rows.Next() {
		var database entities.Database
		err := rows.Scan(&database.ID, &database.ProjectID, &database.Name, &database.Host, &database.Port, &database.Username, &database.Password, &database.Database, &database.SSLMode, &database.CreatedAt, &database.UpdatedAt)
		if err != nil {
			return nil, err
		}
		databases = append(databases, &database)
	}

	return databases, nil
}

func (r *DatabaseRepository) Update(ctx context.Context, database *entities.Database) error {
	query := `
		UPDATE databases
		SET name = $2, host = $3, port = $4, username = $5, password = $6, database = $7, ssl_mode = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, database.ID, database.Name, database.Host, database.Port, database.Username, database.Password, database.Database, database.SSLMode, database.UpdatedAt)
	return err
}

func (r *DatabaseRepository) UpdateTx(ctx context.Context, tx types.Executor, database *entities.Database) error {
	query := `
		UPDATE databases
		SET name = $2, host = $3, port = $4, username = $5, password = $6, database = $7, ssl_mode = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := tx.Exec(ctx, query, database.ID, database.Name, database.Host, database.Port, database.Username, database.Password, database.Database, database.SSLMode, database.UpdatedAt)
	return err
}

func (r *DatabaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM databases WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *DatabaseRepository) DeleteTx(ctx context.Context, tx types.Executor, id uuid.UUID) error {
	query := `DELETE FROM databases WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	return err
}

func (r *DatabaseRepository) List(ctx context.Context, limit, offset int) ([]*entities.Database, error) {
	query := `
		SELECT id, project_id, name, host, port, username, password, database, ssl_mode, created_at, updated_at
		FROM databases 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []*entities.Database
	for rows.Next() {
		var database entities.Database
		err := rows.Scan(&database.ID, &database.ProjectID, &database.Name, &database.Host, &database.Port, &database.Username, &database.Password, &database.Database, &database.SSLMode, &database.CreatedAt, &database.UpdatedAt)
		if err != nil {
			return nil, err
		}
		databases = append(databases, &database)
	}

	return databases, nil
}
