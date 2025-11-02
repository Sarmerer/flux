package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/flow/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TableDataRepository struct {
	db *pgxpool.Pool
}

func NewTableDataRepository(db *pgxpool.Pool) repositories.TableDataRepository {
	return &TableDataRepository{db: db}
}

func (r *TableDataRepository) Query(ctx context.Context, tableName string, limit, offset int) ([]map[string]interface{}, error) {
	sanitizedTable := pgx.Identifier{tableName}.Sanitize()
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT $1 OFFSET $2", sanitizedTable)

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
	}
	defer rows.Close()

	return r.rowsToMaps(rows)
}

func (r *TableDataRepository) Count(ctx context.Context, tableName string) (int64, error) {
	sanitizedTable := pgx.Identifier{tableName}.Sanitize()
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", sanitizedTable)

	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows in table %s: %w", tableName, err)
	}

	return count, nil
}

func (r *TableDataRepository) GetByID(ctx context.Context, tableName string, id interface{}) (map[string]interface{}, error) {
	sanitizedTable := pgx.Identifier{tableName}.Sanitize()
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", sanitizedTable)

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get row from table %s: %w", tableName, err)
	}
	defer rows.Close()

	results, err := r.rowsToMaps(rows)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("row not found")
	}

	return results[0], nil
}

func (r *TableDataRepository) Insert(ctx context.Context, tableName string, data map[string]interface{}) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data provided for insert")
	}

	sanitizedTable := pgx.Identifier{tableName}.Sanitize()

	columns := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	idx := 1

	for col, val := range data {
		columns = append(columns, pgx.Identifier{col}.Sanitize())
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
		values = append(values, val)
		idx++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		sanitizedTable,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	rows, err := r.db.Query(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into table %s: %w", tableName, err)
	}
	defer rows.Close()

	results, err := r.rowsToMaps(rows)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("insert did not return a row")
	}

	return results[0], nil
}

func (r *TableDataRepository) Update(ctx context.Context, tableName string, id interface{}, data map[string]interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided for update")
	}

	sanitizedTable := pgx.Identifier{tableName}.Sanitize()

	setClauses := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data)+1)
	idx := 1

	for col, val := range data {
		if col == "id" {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", pgx.Identifier{col}.Sanitize(), idx))
		values = append(values, val)
		idx++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no valid columns to update")
	}

	values = append(values, id)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		sanitizedTable,
		strings.Join(setClauses, ", "),
		idx,
	)

	cmdTag, err := r.db.Exec(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to update row in table %s: %w", tableName, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("row not found")
	}

	return nil
}

func (r *TableDataRepository) Delete(ctx context.Context, tableName string, id interface{}) error {
	sanitizedTable := pgx.Identifier{tableName}.Sanitize()
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", sanitizedTable)

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete row from table %s: %w", tableName, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("row not found")
	}

	return nil
}

func (r *TableDataRepository) rowsToMaps(rows pgx.Rows) ([]map[string]interface{}, error) {
	fieldDescriptions := rows.FieldDescriptions()
	var results []map[string]interface{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, field := range fieldDescriptions {
			rowMap[string(field.Name)] = values[i]
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
