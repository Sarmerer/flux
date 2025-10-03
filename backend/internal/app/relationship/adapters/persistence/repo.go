package persistence

import (
	"context"
	"fmt"

	"github.com/example/flow/internal/app/relationship/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{ db *pgxpool.Pool }

func NewRepo(db *pgxpool.Pool) *Repo { return &Repo{db: db} }

func (r *Repo) List(tableID int64) ([]domain.Relationship, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT id, source_table_id, target_table_id, kind, created_at, updated_at
        FROM relationships
        WHERE source_table_id=$1 OR target_table_id=$1
        ORDER BY id ASC
    `, tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Relationship
	for rows.Next() {
		var rel domain.Relationship
		if err := rows.Scan(&rel.ID, &rel.SourceTableID, &rel.TargetTableID, &rel.Kind, &rel.CreatedAt, &rel.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, rel)
	}
	return out, rows.Err()
}

func (r *Repo) Create(sourceID, targetID int64, kind domain.Kind) (domain.Relationship, error) {
	var rel domain.Relationship
	err := r.db.QueryRow(context.Background(), `
        INSERT INTO relationships(source_table_id, target_table_id, kind)
        VALUES($1,$2,$3)
        RETURNING id, source_table_id, target_table_id, kind, created_at, updated_at
    `, sourceID, targetID, kind).Scan(&rel.ID, &rel.SourceTableID, &rel.TargetTableID, &rel.Kind, &rel.CreatedAt, &rel.UpdatedAt)
	if err != nil {
		return rel, err
	}

	// DDL management
	switch kind {
	case domain.OneToOne, domain.ManyToOne:
		// Add FK on source pointing to target
		_, err = r.db.Exec(context.Background(), fmt.Sprintf("ALTER TABLE t_%d ADD COLUMN IF NOT EXISTS target_%d BIGINT REFERENCES t_%d(id)", sourceID, targetID, targetID))
	case domain.OneToMany:
		// Add FK on target pointing to source
		_, err = r.db.Exec(context.Background(), fmt.Sprintf("ALTER TABLE t_%d ADD COLUMN IF NOT EXISTS source_%d BIGINT REFERENCES t_%d(id)", targetID, sourceID, sourceID))
	case domain.ManyToMany:
		// Create join table
		_, err = r.db.Exec(context.Background(), fmt.Sprintf("CREATE TABLE IF NOT EXISTS j_%d_%d (a BIGINT REFERENCES t_%d(id), b BIGINT REFERENCES t_%d(id), PRIMARY KEY(a,b))", sourceID, targetID, sourceID, targetID))
	}
	return rel, err
}

func (r *Repo) Delete(id int64) error {
	// For simplicity, leave physical artifacts; advanced cleanup could be added.
	_, err := r.db.Exec(context.Background(), `DELETE FROM relationships WHERE id=$1`, id)
	return err
}
