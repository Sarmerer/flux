package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TableSchemaMutationService handles table schema mutations
type TableSchemaMutationService struct {
	tableRepo   repositories.TableRepository
	dbRepo      repositories.DatabaseRepository
	projectRepo repositories.ProjectRepository
}

// NewTableSchemaMutationService creates a new TableSchemaMutationService
func NewTableSchemaMutationService(
	tableRepo repositories.TableRepository,
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
) *TableSchemaMutationService {
	return &TableSchemaMutationService{
		tableRepo:   tableRepo,
		dbRepo:      dbRepo,
		projectRepo: projectRepo,
	}
}

// ColumnDefinition represents a table column definition
type ColumnDefinition struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Nullable     bool   `json:"nullable"`
	DefaultValue string `json:"default_value,omitempty"`
	PrimaryKey   bool   `json:"primary_key,omitempty"`
	Unique       bool   `json:"unique,omitempty"`
}

// TableSchema represents the complete table schema
type TableSchema struct {
	Columns     []ColumnDefinition     `json:"columns"`
	PrimaryKey  []string               `json:"primary_key,omitempty"`
	Indexes     []IndexDefinition      `json:"indexes,omitempty"`
	ForeignKeys []ForeignKeyDefinition `json:"foreign_keys,omitempty"`
}

// IndexDefinition represents a table index
type IndexDefinition struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// ForeignKeyDefinition represents a foreign key constraint
type ForeignKeyDefinition struct {
	Name             string `json:"name"`
	Column           string `json:"column"`
	ReferencedTable  string `json:"referenced_table"`
	ReferencedColumn string `json:"referenced_column"`
	OnDelete         string `json:"on_delete,omitempty"`
	OnUpdate         string `json:"on_update,omitempty"`
}

// AddColumnRequest represents a request to add a column to a table
type AddColumnRequest struct {
	Column ColumnDefinition `json:"column" validate:"required"`
}

// RemoveColumnRequest represents a request to remove a column from a table
type RemoveColumnRequest struct {
	ColumnName string `json:"column_name" validate:"required"`
}

// ModifyColumnRequest represents a request to modify a column in a table
type ModifyColumnRequest struct {
	ColumnName string           `json:"column_name" validate:"required"`
	Column     ColumnDefinition `json:"column" validate:"required"`
}

// AddForeignKeyRequest represents a request to add a foreign key
type AddForeignKeyRequest struct {
	ForeignKey ForeignKeyDefinition `json:"foreign_key" validate:"required"`
}

// RemoveForeignKeyRequest represents a request to remove a foreign key
type RemoveForeignKeyRequest struct {
	ForeignKeyName string `json:"foreign_key_name" validate:"required"`
}

// CreateTableInDatabase creates a table in the project's PostgreSQL database
func (s *TableSchemaMutationService) CreateTableInDatabase(ctx context.Context, projectID uuid.UUID, tableName string, schema TableSchema) error {
	// Validate table name
	if tableName == "" {
		return fmt.Errorf("table name is required")
	}
	if !isValidTableName(tableName) {
		return fmt.Errorf("invalid table name")
	}

	// Validate schema
	if err := ValidateTableSchema(&schema); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0] // Use first database for now

	// Connect to project database
	pool, err := s.connectToProjectDatabase(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to connect to project database: %w", err)
	}
	defer pool.Close()

	// Generate CREATE TABLE SQL
	createSQL := s.generateCreateTableSQL(tableName, schema)

	// Execute CREATE TABLE
	if _, err := pool.Exec(ctx, createSQL); err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}

	// Create indexes
	for _, index := range schema.Indexes {
		indexSQL := s.generateCreateIndexSQL(tableName, index)
		if _, err := pool.Exec(ctx, indexSQL); err != nil {
			return fmt.Errorf("failed to create index %s: %w", index.Name, err)
		}
	}

	// Create foreign keys
	for _, fk := range schema.ForeignKeys {
		fkSQL := s.generateAddForeignKeySQL(tableName, fk)
		if _, err := pool.Exec(ctx, fkSQL); err != nil {
			return fmt.Errorf("failed to create foreign key %s: %w", fk.Name, err)
		}
	}

	return nil
}

// DropTableFromDatabase drops a table from the project's PostgreSQL database
func (s *TableSchemaMutationService) DropTableFromDatabase(ctx context.Context, projectID uuid.UUID, tableName string) error {
	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	// Connect to project database
	pool, err := s.connectToProjectDatabase(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to connect to project database: %w", err)
	}
	defer pool.Close()

	// Drop table
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)
	if _, err := pool.Exec(ctx, dropSQL); err != nil {
		return fmt.Errorf("failed to drop table %s: %w", tableName, err)
	}

	return nil
}

// AddColumnToTable adds a column to an existing table
func (s *TableSchemaMutationService) AddColumnToTable(ctx context.Context, projectID uuid.UUID, tableName string, req *AddColumnRequest) error {
	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	// Connect to project database
	pool, err := s.connectToProjectDatabase(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to connect to project database: %w", err)
	}
	defer pool.Close()

	// Generate ALTER TABLE ADD COLUMN SQL
	addColumnSQL := s.generateAddColumnSQL(tableName, req.Column)

	// Execute ALTER TABLE
	if _, err := pool.Exec(ctx, addColumnSQL); err != nil {
		return fmt.Errorf("failed to add column %s to table %s: %w", req.Column.Name, tableName, err)
	}

	return nil
}

// RemoveColumnFromTable removes a column from an existing table
func (s *TableSchemaMutationService) RemoveColumnFromTable(ctx context.Context, projectID uuid.UUID, tableName string, req *RemoveColumnRequest) error {
	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	// Connect to project database
	pool, err := s.connectToProjectDatabase(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to connect to project database: %w", err)
	}
	defer pool.Close()

	// Generate ALTER TABLE DROP COLUMN SQL
	dropColumnSQL := fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s", tableName, req.ColumnName)

	// Execute ALTER TABLE
	if _, err := pool.Exec(ctx, dropColumnSQL); err != nil {
		return fmt.Errorf("failed to remove column %s from table %s: %w", req.ColumnName, tableName, err)
	}

	return nil
}

// ModifyColumnInTable modifies a column in an existing table
func (s *TableSchemaMutationService) ModifyColumnInTable(ctx context.Context, projectID uuid.UUID, tableName string, req *ModifyColumnRequest) error {
	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	// Connect to project database
	pool, err := s.connectToProjectDatabase(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to connect to project database: %w", err)
	}
	defer pool.Close()

	// Generate ALTER TABLE ALTER COLUMN SQL
	modifyColumnSQL := s.generateModifyColumnSQL(tableName, req.ColumnName, req.Column)

	// Execute ALTER TABLE
	if _, err := pool.Exec(ctx, modifyColumnSQL); err != nil {
		return fmt.Errorf("failed to modify column %s in table %s: %w", req.ColumnName, tableName, err)
	}

	return nil
}

// AddForeignKeyToTable adds a foreign key constraint to a table
func (s *TableSchemaMutationService) AddForeignKeyToTable(ctx context.Context, projectID uuid.UUID, tableName string, req *AddForeignKeyRequest) error {
	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	// Connect to project database
	pool, err := s.connectToProjectDatabase(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to connect to project database: %w", err)
	}
	defer pool.Close()

	// Generate ADD FOREIGN KEY SQL
	fkSQL := s.generateAddForeignKeySQL(tableName, req.ForeignKey)

	// Execute ALTER TABLE
	if _, err := pool.Exec(ctx, fkSQL); err != nil {
		return fmt.Errorf("failed to add foreign key %s to table %s: %w", req.ForeignKey.Name, tableName, err)
	}

	return nil
}

// RemoveForeignKeyFromTable removes a foreign key constraint from a table
func (s *TableSchemaMutationService) RemoveForeignKeyFromTable(ctx context.Context, projectID uuid.UUID, tableName string, req *RemoveForeignKeyRequest) error {
	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("no database found for project: %w", err)
	}

	database := databases[0]

	// Connect to project database
	pool, err := s.connectToProjectDatabase(ctx, database)
	if err != nil {
		return fmt.Errorf("failed to connect to project database: %w", err)
	}
	defer pool.Close()

	// Generate DROP FOREIGN KEY SQL
	dropFKSQL := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s", tableName, req.ForeignKeyName)

	// Execute ALTER TABLE
	if _, err := pool.Exec(ctx, dropFKSQL); err != nil {
		return fmt.Errorf("failed to remove foreign key %s from table %s: %w", req.ForeignKeyName, tableName, err)
	}

	return nil
}

// connectToProjectDatabase creates a connection to a project's database
func (s *TableSchemaMutationService) connectToProjectDatabase(ctx context.Context, database *entities.Database) (*pgxpool.Pool, error) {
	dsn := database.GetConnectionString()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return pool, nil
}

// generateCreateTableSQL generates CREATE TABLE SQL from schema
func (s *TableSchemaMutationService) generateCreateTableSQL(tableName string, schema TableSchema) string {
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

	// Add primary key constraint if specified
	if len(schema.PrimaryKey) > 0 {
		columns = append(columns, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(schema.PrimaryKey, ", ")))
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))
}

// generateAddColumnSQL generates ALTER TABLE ADD COLUMN SQL
func (s *TableSchemaMutationService) generateAddColumnSQL(tableName string, column ColumnDefinition) string {
	colDef := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, column.Name, column.Type)

	if !column.Nullable {
		colDef += " NOT NULL"
	}

	if column.DefaultValue != "" {
		colDef += fmt.Sprintf(" DEFAULT %s", column.DefaultValue)
	}

	return colDef
}

// generateModifyColumnSQL generates ALTER TABLE ALTER COLUMN SQL
func (s *TableSchemaMutationService) generateModifyColumnSQL(tableName, oldColumnName string, newColumn ColumnDefinition) string {
	var modifications []string

	// Change column type
	modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s TYPE %s", oldColumnName, newColumn.Type))

	// Change nullability
	if !newColumn.Nullable {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s SET NOT NULL", oldColumnName))
	} else {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s DROP NOT NULL", oldColumnName))
	}

	// Change default value
	if newColumn.DefaultValue != "" {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s SET DEFAULT %s", oldColumnName, newColumn.DefaultValue))
	} else {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s DROP DEFAULT", oldColumnName))
	}

	// Rename column if name changed
	if oldColumnName != newColumn.Name {
		modifications = append(modifications, fmt.Sprintf("RENAME COLUMN %s TO %s", oldColumnName, newColumn.Name))
	}

	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(modifications, ", "))
}

// generateCreateIndexSQL generates CREATE INDEX SQL
func (s *TableSchemaMutationService) generateCreateIndexSQL(tableName string, index IndexDefinition) string {
	unique := ""
	if index.Unique {
		unique = "UNIQUE "
	}

	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)",
		unique, index.Name, tableName, strings.Join(index.Columns, ", "))
}

// generateAddForeignKeySQL generates ADD FOREIGN KEY SQL
func (s *TableSchemaMutationService) generateAddForeignKeySQL(tableName string, fk ForeignKeyDefinition) string {
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
