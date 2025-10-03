package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/example/flow/internal/app/column/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{ db *pgxpool.Pool }

func NewRepo(db *pgxpool.Pool) *Repo { return &Repo{db: db} }

func (r *Repo) List(tableID int64) ([]domain.Column, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT id, table_id, name, type, required, unique_col, default_value, created_at, updated_at
        FROM columns WHERE table_id=$1 ORDER BY id ASC`, tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Column
	for rows.Next() {
		var c domain.Column
		if err := rows.Scan(&c.ID, &c.TableID, &c.Name, &c.Type, &c.Required, &c.Unique, &c.DefaultValue, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repo) Create(tableID int64, name string, typ domain.ColumnType, required bool, unique bool, defaultValue *string) (domain.Column, error) {
	var c domain.Column
	err := r.db.QueryRow(context.Background(), `
        INSERT INTO columns(table_id,name,type,required,unique_col,default_value)
        VALUES($1,$2,$3,$4,$5,$6)
        RETURNING id, table_id, name, type, required, unique_col, default_value, created_at, updated_at
    `, tableID, name, typ, required, unique, defaultValue).Scan(&c.ID, &c.TableID, &c.Name, &c.Type, &c.Required, &c.Unique, &c.DefaultValue, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return domain.Column{}, err
	}
	// Apply physical ALTER TABLE
	safe := safeIdent(name)
	tableName := physicalTableName(tableID)
	var ddl string
	switch typ {
	case domain.ColumnText:
		ddl = fmt.Sprintf("ALTER TABLE %s ADD COLUMN \"%s\" TEXT", tableName, safe)
	case domain.ColumnInteger:
		ddl = fmt.Sprintf("ALTER TABLE %s ADD COLUMN \"%s\" BIGINT", tableName, safe)
	case domain.ColumnBoolean:
		ddl = fmt.Sprintf("ALTER TABLE %s ADD COLUMN \"%s\" BOOLEAN", tableName, safe)
	case domain.ColumnTimestamp:
		ddl = fmt.Sprintf("ALTER TABLE %s ADD COLUMN \"%s\" TIMESTAMPTZ", tableName, safe)
	}
	if required {
		ddl += " NOT NULL"
	}
	if defaultValue != nil && *defaultValue != "" {
		ddl += fmt.Sprintf(" DEFAULT '%s'", strings.ReplaceAll(*defaultValue, "'", "''"))
	}
	if _, err := r.db.Exec(context.Background(), ddl); err != nil {
		return c, err
	}
	if unique {
		if _, err := r.db.Exec(context.Background(), fmt.Sprintf("CREATE UNIQUE INDEX ON %s (\"%s\")", tableName, safe)); err != nil {
			return c, err
		}
	}
	return c, nil
}

func (r *Repo) Update(id int64, name string, required bool, unique bool, defaultValue *string) (domain.Column, error) {
	var c domain.Column
	err := r.db.QueryRow(context.Background(), `
        UPDATE columns SET name=$1, required=$2, unique_col=$3, default_value=$4, updated_at=NOW()
        WHERE id=$5 RETURNING id, table_id, name, type, required, unique_col, default_value, created_at, updated_at
    `, name, required, unique, defaultValue, id).Scan(&c.ID, &c.TableID, &c.Name, &c.Type, &c.Required, &c.Unique, &c.DefaultValue, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repo) Delete(id int64) error {
	// Drop physical column as well
	var tableID int64
	var name string
	if err := r.db.QueryRow(context.Background(), `SELECT table_id, name FROM columns WHERE id=$1`, id).Scan(&tableID, &name); err != nil {
		return err
	}
	if _, err := r.db.Exec(context.Background(), `DELETE FROM columns WHERE id=$1`, id); err != nil {
		return err
	}
	_, err := r.db.Exec(context.Background(), fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS \"%s\"", physicalTableName(tableID), safeIdent(name)))
	return err
}

func physicalTableName(tableID int64) string {
	return fmt.Sprintf("t_%d", tableID)
}

func safeIdent(s string) string { return strings.ReplaceAll(s, "\"", "") }
