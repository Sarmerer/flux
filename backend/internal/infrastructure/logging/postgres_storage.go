package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresLogStorage struct {
	db *pgxpool.Pool
}

func NewPostgresLogStorage(db *pgxpool.Pool) *PostgresLogStorage {
	return &PostgresLogStorage{
		db: db,
	}
}

func (s *PostgresLogStorage) Store(ctx context.Context, entry *LogEntry) error {

	fieldsJSON, err := json.Marshal(entry.Fields)
	if err != nil {
		return fmt.Errorf("failed to marshal fields: %w", err)
	}

	query := `
		INSERT INTO logs (
			id, level, message, timestamp,
			project_id, database_id, table_id, workflow_id, user_id,
			request_id, operation, fields
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	_, err = s.db.Exec(ctx, query,
		entry.ID,
		entry.Level,
		entry.Message,
		entry.Timestamp,
		entry.ProjectID,
		entry.DatabaseID,
		entry.TableID,
		entry.WorkflowID,
		entry.UserID,
		entry.RequestID,
		entry.Operation,
		fieldsJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to store log entry: %w", err)
	}

	return nil
}

func (s *PostgresLogStorage) Query(ctx context.Context, filter *LogFilter) ([]*LogEntry, error) {
	if filter == nil {
		filter = DefaultFilter()
	}

	query := `SELECT id, level, message, timestamp, project_id, database_id, table_id, workflow_id, user_id, request_id, operation, fields FROM logs WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if filter.ProjectID != nil {
		query += fmt.Sprintf(" AND project_id = $%d", argPos)
		args = append(args, filter.ProjectID)
		argPos++
	}

	if filter.DatabaseID != nil {
		query += fmt.Sprintf(" AND database_id = $%d", argPos)
		args = append(args, filter.DatabaseID)
		argPos++
	}

	if filter.TableID != nil {
		query += fmt.Sprintf(" AND table_id = $%d", argPos)
		args = append(args, filter.TableID)
		argPos++
	}

	if filter.WorkflowID != nil {
		query += fmt.Sprintf(" AND workflow_id = $%d", argPos)
		args = append(args, filter.WorkflowID)
		argPos++
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, filter.UserID)
		argPos++
	}

	if filter.Level != nil {
		query += fmt.Sprintf(" AND level = $%d", argPos)
		args = append(args, filter.Level)
		argPos++
	}

	if filter.Operation != nil {
		query += fmt.Sprintf(" AND operation = $%d", argPos)
		args = append(args, filter.Operation)
		argPos++
	}

	if filter.StartTime != nil {
		query += fmt.Sprintf(" AND timestamp >= $%d", argPos)
		args = append(args, filter.StartTime)
		argPos++
	}

	if filter.EndTime != nil {
		query += fmt.Sprintf(" AND timestamp <= $%d", argPos)
		args = append(args, filter.EndTime)
		argPos++
	}

	if filter.OrderBy == "timestamp_asc" {
		query += " ORDER BY timestamp ASC"
	} else {
		query += " ORDER BY timestamp DESC"
	}

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filter.Limit)
		argPos++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, filter.Offset)
		argPos++
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	entries := make([]*LogEntry, 0)
	for rows.Next() {
		entry := &LogEntry{}
		var fieldsJSON []byte

		err := rows.Scan(
			&entry.ID,
			&entry.Level,
			&entry.Message,
			&entry.Timestamp,
			&entry.ProjectID,
			&entry.DatabaseID,
			&entry.TableID,
			&entry.WorkflowID,
			&entry.UserID,
			&entry.RequestID,
			&entry.Operation,
			&fieldsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan log entry: %w", err)
		}

		if len(fieldsJSON) > 0 {
			if err := json.Unmarshal(fieldsJSON, &entry.Fields); err != nil {
				return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
			}
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating log entries: %w", err)
	}

	return entries, nil
}

func (s *PostgresLogStorage) Count(ctx context.Context, filter *LogFilter) (int64, error) {
	if filter == nil {
		filter = DefaultFilter()
	}

	query := `SELECT COUNT(*) FROM logs WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if filter.ProjectID != nil {
		query += fmt.Sprintf(" AND project_id = $%d", argPos)
		args = append(args, filter.ProjectID)
		argPos++
	}

	if filter.DatabaseID != nil {
		query += fmt.Sprintf(" AND database_id = $%d", argPos)
		args = append(args, filter.DatabaseID)
		argPos++
	}

	if filter.TableID != nil {
		query += fmt.Sprintf(" AND table_id = $%d", argPos)
		args = append(args, filter.TableID)
		argPos++
	}

	if filter.WorkflowID != nil {
		query += fmt.Sprintf(" AND workflow_id = $%d", argPos)
		args = append(args, filter.WorkflowID)
		argPos++
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, filter.UserID)
		argPos++
	}

	if filter.Level != nil {
		query += fmt.Sprintf(" AND level = $%d", argPos)
		args = append(args, filter.Level)
		argPos++
	}

	if filter.Operation != nil {
		query += fmt.Sprintf(" AND operation = $%d", argPos)
		args = append(args, filter.Operation)
		argPos++
	}

	if filter.StartTime != nil {
		query += fmt.Sprintf(" AND timestamp >= $%d", argPos)
		args = append(args, filter.StartTime)
		argPos++
	}

	if filter.EndTime != nil {
		query += fmt.Sprintf(" AND timestamp <= $%d", argPos)
		args = append(args, filter.EndTime)
		argPos++
	}

	var count int64
	err := s.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count logs: %w", err)
	}

	return count, nil
}

func (s *PostgresLogStorage) DeleteOlderThan(ctx context.Context, duration time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-duration)

	query := `DELETE FROM logs WHERE timestamp < $1`

	result, err := s.db.Exec(ctx, query, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old logs: %w", err)
	}

	return result.RowsAffected(), nil
}

func MigrateLogsTable(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS logs (
		id UUID PRIMARY KEY,
		level VARCHAR(10) NOT NULL,
		message TEXT NOT NULL,
		timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
		project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
		database_id UUID,
		table_id UUID,
		workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE,
		user_id UUID REFERENCES users(id) ON DELETE SET NULL,
		request_id VARCHAR(255),
		operation VARCHAR(255),
		fields JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Indexes for efficient querying
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_logs_project_id ON logs(project_id);
	CREATE INDEX IF NOT EXISTS idx_logs_database_id ON logs(database_id);
	CREATE INDEX IF NOT EXISTS idx_logs_table_id ON logs(table_id);
	CREATE INDEX IF NOT EXISTS idx_logs_workflow_id ON logs(workflow_id);
	CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);
	CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
	CREATE INDEX IF NOT EXISTS idx_logs_operation ON logs(operation);

	-- Composite index for common queries
	CREATE INDEX IF NOT EXISTS idx_logs_project_timestamp ON logs(project_id, timestamp DESC);
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create logs table: %w", err)
	}

	return nil
}
