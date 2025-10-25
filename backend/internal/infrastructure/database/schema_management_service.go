package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/flow/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ColumnDefinition struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Nullable     bool   `json:"nullable"`
	DefaultValue string `json:"default_value,omitempty"`
	PrimaryKey   bool   `json:"primary_key,omitempty"`
	Unique       bool   `json:"unique,omitempty"`
}

type TableSchema struct {
	Columns     []ColumnDefinition     `json:"columns"`
	PrimaryKey  []string               `json:"primary_key,omitempty"`
	Indexes     []IndexDefinition      `json:"indexes,omitempty"`
	ForeignKeys []ForeignKeyDefinition `json:"foreign_keys,omitempty"`
}

type IndexDefinition struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

type ForeignKeyDefinition struct {
	Name             string `json:"name"`
	Column           string `json:"column"`
	ReferencedTable  string `json:"referenced_table"`
	ReferencedColumn string `json:"referenced_column"`
	OnDelete         string `json:"on_delete,omitempty"`
	OnUpdate         string `json:"on_update,omitempty"`
}

type SchemaManagementService struct {
	connService *ConnectionService
}

func NewSchemaManagementService(connService *ConnectionService) *SchemaManagementService {
	return &SchemaManagementService{
		connService: connService,
	}
}

func (s *SchemaManagementService) CreateTable(ctx context.Context, database *entities.Database, tableName string, schema TableSchema) error {
	createSQL := s.buildCreateTableSQL(tableName, schema)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, createSQL); err != nil {
			return fmt.Errorf("failed to create table %s: %w", tableName, err)
		}

		for _, index := range schema.Indexes {
			indexSQL := s.buildCreateIndexSQL(tableName, index)
			if _, err := pool.Exec(ctx, indexSQL); err != nil {
				return fmt.Errorf("failed to create index %s: %w", index.Name, err)
			}
		}

		for _, fk := range schema.ForeignKeys {
			fkSQL := s.buildAddForeignKeySQL(tableName, fk)
			if _, err := pool.Exec(ctx, fkSQL); err != nil {
				return fmt.Errorf("failed to create foreign key %s: %w", fk.Name, err)
			}
		}

		return nil
	})
}

func (s *SchemaManagementService) DropTable(ctx context.Context, database *entities.Database, tableName string) error {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, dropSQL); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", tableName, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) AddColumn(ctx context.Context, database *entities.Database, tableName string, column ColumnDefinition) error {
	addColumnSQL := s.buildAddColumnSQL(tableName, column)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, addColumnSQL); err != nil {
			return fmt.Errorf("failed to add column %s to table %s: %w", column.Name, tableName, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) RemoveColumn(ctx context.Context, database *entities.Database, tableName string, columnName string) error {
	dropColumnSQL := fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s", tableName, columnName)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, dropColumnSQL); err != nil {
			return fmt.Errorf("failed to remove column %s from table %s: %w", columnName, tableName, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) ModifyColumn(ctx context.Context, database *entities.Database, tableName string, oldColumnName string, newColumn ColumnDefinition) error {
	modifyColumnSQL := s.buildModifyColumnSQL(tableName, oldColumnName, newColumn)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, modifyColumnSQL); err != nil {
			return fmt.Errorf("failed to modify column %s in table %s: %w", oldColumnName, tableName, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) AddForeignKey(ctx context.Context, database *entities.Database, tableName string, fk ForeignKeyDefinition) error {
	fkSQL := s.buildAddForeignKeySQL(tableName, fk)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, fkSQL); err != nil {
			return fmt.Errorf("failed to add foreign key %s to table %s: %w", fk.Name, tableName, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) RemoveForeignKey(ctx context.Context, database *entities.Database, tableName string, foreignKeyName string) error {
	dropFKSQL := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s", tableName, foreignKeyName)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, dropFKSQL); err != nil {
			return fmt.Errorf("failed to remove foreign key %s from table %s: %w", foreignKeyName, tableName, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) buildCreateTableSQL(tableName string, schema TableSchema) string {
	var columns []string

	for _, col := range schema.Columns {
		colDef := fmt.Sprintf("%s %s", col.Name, col.Type)

		if !col.Nullable {
			colDef += " NOT NULL"
		}

		if col.DefaultValue != "" {
			colDef += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
		}

		if col.Unique {
			colDef += " UNIQUE"
		}

		columns = append(columns, colDef)
	}

	if len(schema.PrimaryKey) > 0 {
		columns = append(columns, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(schema.PrimaryKey, ", ")))
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))
}

func (s *SchemaManagementService) buildAddColumnSQL(tableName string, column ColumnDefinition) string {
	colDef := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, column.Name, column.Type)

	if !column.Nullable {
		colDef += " NOT NULL"
	}

	if column.DefaultValue != "" {
		colDef += fmt.Sprintf(" DEFAULT %s", column.DefaultValue)
	}

	return colDef
}

func (s *SchemaManagementService) buildModifyColumnSQL(tableName, oldColumnName string, newColumn ColumnDefinition) string {
	var modifications []string

	modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s TYPE %s", oldColumnName, newColumn.Type))

	if !newColumn.Nullable {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s SET NOT NULL", oldColumnName))
	} else {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s DROP NOT NULL", oldColumnName))
	}

	if newColumn.DefaultValue != "" {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s SET DEFAULT %s", oldColumnName, newColumn.DefaultValue))
	} else {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s DROP DEFAULT", oldColumnName))
	}

	if oldColumnName != newColumn.Name {
		modifications = append(modifications, fmt.Sprintf("RENAME COLUMN %s TO %s", oldColumnName, newColumn.Name))
	}

	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(modifications, ", "))
}

func (s *SchemaManagementService) buildCreateIndexSQL(tableName string, index IndexDefinition) string {
	unique := ""
	if index.Unique {
		unique = "UNIQUE "
	}

	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)",
		unique, index.Name, tableName, strings.Join(index.Columns, ", "))
}

func (s *SchemaManagementService) buildAddForeignKeySQL(tableName string, fk ForeignKeyDefinition) string {
	sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
		tableName, fk.Name, fk.Column, fk.ReferencedTable, fk.ReferencedColumn)

	if fk.OnDelete != "" {
		sql += fmt.Sprintf(" ON DELETE %s", fk.OnDelete)
	}

	if fk.OnUpdate != "" {
		sql += fmt.Sprintf(" ON UPDATE %s", fk.OnUpdate)
	}

	return sql
}
