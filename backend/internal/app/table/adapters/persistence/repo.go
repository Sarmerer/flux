package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/example/flow/internal/app/table/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo { return &Repo{db: db} }

func (r *Repo) ListTables() ([]domain.Table, error) {
	rows, err := r.db.Query(context.Background(), `SELECT id, name, created_at, updated_at FROM tables ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Table
	for rows.Next() {
		var t domain.Table
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) CreateTable(name string) (domain.Table, error) {
	var t domain.Table
	now := time.Now().UTC()
	err := r.db.QueryRow(context.Background(), `INSERT INTO tables(name, created_at, updated_at) VALUES($1,$2,$3) RETURNING id, name, created_at, updated_at`, name, now, now).
		Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return t, err
	}
	_, err = r.db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS t_"+fmt.Sprint(t.ID)+" (id BIGSERIAL PRIMARY KEY)")
	return t, err
}

func (r *Repo) UpdateTable(id int64, name string) (domain.Table, error) {
	var t domain.Table
	err := r.db.QueryRow(context.Background(), `UPDATE tables SET name=$1, updated_at=NOW() AT TIME ZONE 'UTC' WHERE id=$2 RETURNING id, name, created_at, updated_at`, name, id).
		Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *Repo) DeleteTable(id int64) error {
	ct, err := r.db.Exec(context.Background(), `DELETE FROM tables WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}
