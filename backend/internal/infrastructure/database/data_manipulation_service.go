package database

import (
	"context"
	"fmt"

	"github.com/flow/internal/domain/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DataManipulationService struct {
	connService *ConnectionService
}

func NewDataManipulationService(connService *ConnectionService) *DataManipulationService {
	return &DataManipulationService{
		connService: connService,
	}
}

func (s *DataManipulationService) UpdateRow(ctx context.Context, database *entities.Database, tableName string, rowID string, updates map[string]any) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	setClause := ""
	args := []any{}
	argPos := 1

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", pgx.Identifier{key}.Sanitize(), argPos)
		args = append(args, value)
		argPos++
	}

	args = append(args, rowID)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", pgx.Identifier{tableName}.Sanitize(), setClause, argPos)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, query, args...); err != nil {
			return fmt.Errorf("Failed to update row in table %s: %w", tableName, err)
		}
		return nil
	})
}

func (s *DataManipulationService) CreateRow(ctx context.Context, database *entities.Database, tableName string, data map[string]any) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	columns := ""
	placeholders := ""
	args := []any{}
	argPos := 1

	for key, value := range data {
		if columns != "" {
			columns += ", "
			placeholders += ", "
		}
		columns += pgx.Identifier{key}.Sanitize()
		placeholders += fmt.Sprintf("$%d", argPos)
		args = append(args, value)
		argPos++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", pgx.Identifier{tableName}.Sanitize(), columns, placeholders)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, query, args...); err != nil {
			return fmt.Errorf("Failed to insert row into table %s: %w", tableName, err)
		}
		return nil
	})
}

func (s *DataManipulationService) DeleteRow(ctx context.Context, database *entities.Database, tableName string, rowID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", pgx.Identifier{tableName}.Sanitize())

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, query, rowID); err != nil {
			return fmt.Errorf("Failed to delete row from table %s: %w", tableName, err)
		}
		return nil
	})
}
