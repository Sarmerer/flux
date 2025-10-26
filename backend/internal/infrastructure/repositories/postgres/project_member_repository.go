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

type ProjectMemberRepository struct {
	db *pgxpool.Pool
}

func NewProjectMemberRepository(db *pgxpool.Pool) repositories.ProjectMemberRepository {
	return &ProjectMemberRepository{db: db}
}

func (r *ProjectMemberRepository) Create(ctx context.Context, member *entities.ProjectMember) error {
	query := `
		INSERT INTO project_members (id, project_id, user_id, role, permissions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, member.ID, member.ProjectID, member.UserID, member.Role, member.Permissions, member.CreatedAt, member.UpdatedAt)
	return err
}

func (r *ProjectMemberRepository) CreateTx(ctx context.Context, tx types.Executor, member *entities.ProjectMember) error {
	query := `
		INSERT INTO project_members (id, project_id, user_id, role, permissions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(ctx, query, member.ID, member.ProjectID, member.UserID, member.Role, member.Permissions, member.CreatedAt, member.UpdatedAt)
	return err
}

func (r *ProjectMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, role, permissions, created_at, updated_at
		FROM project_members WHERE id = $1
	`
	var member entities.ProjectMember
	err := r.db.QueryRow(ctx, query, id).Scan(
		&member.ID, &member.ProjectID, &member.UserID, &member.Role, &member.Permissions, &member.CreatedAt, &member.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project member not found")
		}
		return nil, err
	}
	return &member, nil
}

func (r *ProjectMemberRepository) GetByProjectAndUser(ctx context.Context, projectID, userID uuid.UUID) (*entities.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, role, permissions, created_at, updated_at
		FROM project_members WHERE project_id = $1 AND user_id = $2
	`
	var member entities.ProjectMember
	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(
		&member.ID, &member.ProjectID, &member.UserID, &member.Role, &member.Permissions, &member.CreatedAt, &member.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project member not found")
		}
		return nil, err
	}
	return &member, nil
}

func (r *ProjectMemberRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, role, permissions, created_at, updated_at
		FROM project_members WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*entities.ProjectMember
	for rows.Next() {
		var member entities.ProjectMember
		err := rows.Scan(&member.ID, &member.ProjectID, &member.UserID, &member.Role, &member.Permissions, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	return members, nil
}

func (r *ProjectMemberRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, role, permissions, created_at, updated_at
		FROM project_members WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*entities.ProjectMember
	for rows.Next() {
		var member entities.ProjectMember
		err := rows.Scan(&member.ID, &member.ProjectID, &member.UserID, &member.Role, &member.Permissions, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	return members, nil
}

func (r *ProjectMemberRepository) Update(ctx context.Context, member *entities.ProjectMember) error {
	query := `
		UPDATE project_members
		SET role = $2, permissions = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, member.ID, member.Role, member.Permissions, member.UpdatedAt)
	return err
}

func (r *ProjectMemberRepository) UpdateTx(ctx context.Context, tx types.Executor, member *entities.ProjectMember) error {
	query := `
		UPDATE project_members
		SET role = $2, permissions = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := tx.Exec(ctx, query, member.ID, member.Role, member.Permissions, member.UpdatedAt)
	return err
}

func (r *ProjectMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM project_members WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ProjectMemberRepository) DeleteTx(ctx context.Context, tx types.Executor, id uuid.UUID) error {
	query := `DELETE FROM project_members WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	return err
}
