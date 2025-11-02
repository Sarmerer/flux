package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/flow/internal/domain/entities"
	"github.com/jackc/pgx/v5"
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

func (s *SchemaManagementService) CreateTable(ctx context.Context, database *entities.Database, table *entities.Table, schema TableSchema) error {
	createSQL := s.buildCreateTableSQL(table.Name, schema)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, createSQL); err != nil {
			return fmt.Errorf("Failed to create table %s: %w", table.Name, err)
		}

		for _, index := range schema.Indexes {
			indexSQL := s.buildCreateIndexSQL(table.Name, index)
			if _, err := pool.Exec(ctx, indexSQL); err != nil {
				return fmt.Errorf("Failed to create index %s: %w", index.Name, err)
			}
		}

		for _, fk := range schema.ForeignKeys {
			fkSQL := s.buildAddForeignKeySQL(table.Name, fk)
			if _, err := pool.Exec(ctx, fkSQL); err != nil {
				return fmt.Errorf("Failed to create foreign key %s: %w", fk.Name, err)
			}
		}

		return nil
	})
}

func (s *SchemaManagementService) DropTable(ctx context.Context, database *entities.Database, table *entities.Table) error {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", pgx.Identifier{table.Name}.Sanitize())

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, dropSQL); err != nil {
			return fmt.Errorf("Failed to drop table %s: %w", table.Name, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) AddColumn(ctx context.Context, database *entities.Database, table *entities.Table, column ColumnDefinition) error {
	addColumnSQL := s.buildAddColumnSQL(table.Name, column)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, addColumnSQL); err != nil {
			return fmt.Errorf("Failed to add column %s to table %s: %w", column.Name, table.Name, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) RemoveColumn(ctx context.Context, database *entities.Database, table *entities.Table, columnName string) error {
	dropColumnSQL := fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s",
		pgx.Identifier{table.Name}.Sanitize(),
		pgx.Identifier{columnName}.Sanitize())

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, dropColumnSQL); err != nil {
			return fmt.Errorf("Failed to remove column %s from table %s: %w", columnName, table.Name, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) ModifyColumn(ctx context.Context, database *entities.Database, table *entities.Table, oldColumnName string, newColumn ColumnDefinition) error {
	modifyColumnSQL := s.buildModifyColumnSQL(table.Name, oldColumnName, newColumn)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, modifyColumnSQL); err != nil {
			return fmt.Errorf("Failed to modify column %s in table %s: %w", oldColumnName, table.Name, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) AddForeignKey(ctx context.Context, database *entities.Database, table *entities.Table, fk ForeignKeyDefinition) error {
	fkSQL := s.buildAddForeignKeySQL(table.Name, fk)

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, fkSQL); err != nil {
			return fmt.Errorf("Failed to add foreign key %s to table %s: %w", fk.Name, table.Name, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) RemoveForeignKey(ctx context.Context, database *entities.Database, table *entities.Table, foreignKeyName string) error {
	dropFKSQL := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s",
		pgx.Identifier{table.Name}.Sanitize(),
		pgx.Identifier{foreignKeyName}.Sanitize())

	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		if _, err := pool.Exec(ctx, dropFKSQL); err != nil {
			return fmt.Errorf("Failed to remove foreign key %s from table %s: %w", foreignKeyName, table.Name, err)
		}
		return nil
	})
}

func (s *SchemaManagementService) UpdateTableSchema(ctx context.Context, database *entities.Database, table *entities.Table, schema TableSchema) error {
	return s.connService.ExecuteWithProjectDB(ctx, database, func(pool *pgxpool.Pool) error {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		tempTableName := fmt.Sprintf("%s_new", table.Name)
		createSQL := s.buildCreateTableSQL(tempTableName, schema)

		if _, err := tx.Exec(ctx, createSQL); err != nil {
			return fmt.Errorf("failed to create temporary table: %w", err)
		}

		for _, index := range schema.Indexes {
			indexSQL := s.buildCreateIndexSQL(tempTableName, index)
			if _, err := tx.Exec(ctx, indexSQL); err != nil {
				return fmt.Errorf("failed to create index %s: %w", index.Name, err)
			}
		}

		columnNames := make([]string, len(schema.Columns))
		for i, col := range schema.Columns {
			columnNames[i] = pgx.Identifier{col.Name}.Sanitize()
		}
		columnsStr := strings.Join(columnNames, ", ")

		copySQL := fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s",
			pgx.Identifier{tempTableName}.Sanitize(),
			columnsStr,
			columnsStr,
			pgx.Identifier{table.Name}.Sanitize())

		if _, err := tx.Exec(ctx, copySQL); err != nil {
			return fmt.Errorf("failed to copy data to temporary table: %w", err)
		}

		dropSQL := fmt.Sprintf("DROP TABLE %s CASCADE", pgx.Identifier{table.Name}.Sanitize())
		if _, err := tx.Exec(ctx, dropSQL); err != nil {
			return fmt.Errorf("failed to drop original table: %w", err)
		}

		renameSQL := fmt.Sprintf("ALTER TABLE %s RENAME TO %s",
			pgx.Identifier{tempTableName}.Sanitize(),
			pgx.Identifier{table.Name}.Sanitize())
		if _, err := tx.Exec(ctx, renameSQL); err != nil {
			return fmt.Errorf("failed to rename temporary table: %w", err)
		}

		for _, fk := range schema.ForeignKeys {
			fkSQL := s.buildAddForeignKeySQL(table.Name, fk)
			if _, err := tx.Exec(ctx, fkSQL); err != nil {
				return fmt.Errorf("failed to create foreign key %s: %w", fk.Name, err)
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	})
}

func (s *SchemaManagementService) buildCreateTableSQL(tableName string, schema TableSchema) string {
	var columns []string

	for _, col := range schema.Columns {
		colName := pgx.Identifier{col.Name}.Sanitize()
		colDef := fmt.Sprintf("%s %s", colName, col.Type)

		if !col.Nullable {
			colDef += " NOT NULL"
		}

		if col.DefaultValue != "" {
			colDef += " DEFAULT " + col.DefaultValue
		}

		if col.Unique {
			colDef += " UNIQUE"
		}

		columns = append(columns, colDef)
	}

	if len(schema.PrimaryKey) > 0 {
		pkColumns := make([]string, len(schema.PrimaryKey))
		for i, col := range schema.PrimaryKey {
			pkColumns[i] = pgx.Identifier{col}.Sanitize()
		}
		columns = append(columns, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pkColumns, ", ")))
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", pgx.Identifier{tableName}.Sanitize(), strings.Join(columns, ", "))
}

func (s *SchemaManagementService) buildAddColumnSQL(tableName string, column ColumnDefinition) string {
	colDef := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s",
		pgx.Identifier{tableName}.Sanitize(),
		pgx.Identifier{column.Name}.Sanitize(),
		column.Type)

	if !column.Nullable {
		colDef += " NOT NULL"
	}

	if column.DefaultValue != "" {
		colDef += " DEFAULT " + column.DefaultValue
	}

	return colDef
}

func (s *SchemaManagementService) buildModifyColumnSQL(tableName, oldColumnName string, newColumn ColumnDefinition) string {
	var modifications []string
	oldColSanitized := pgx.Identifier{oldColumnName}.Sanitize()

	modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s TYPE %s", oldColSanitized, newColumn.Type))

	if !newColumn.Nullable {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s SET NOT NULL", oldColSanitized))
	} else {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s DROP NOT NULL", oldColSanitized))
	}

	if newColumn.DefaultValue != "" {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s SET DEFAULT %s", oldColSanitized, newColumn.DefaultValue))
	} else {
		modifications = append(modifications, fmt.Sprintf("ALTER COLUMN %s DROP DEFAULT", oldColSanitized))
	}

	if oldColumnName != newColumn.Name {
		modifications = append(modifications, fmt.Sprintf("RENAME COLUMN %s TO %s",
			oldColSanitized,
			pgx.Identifier{newColumn.Name}.Sanitize()))
	}

	return fmt.Sprintf("ALTER TABLE %s %s", pgx.Identifier{tableName}.Sanitize(), strings.Join(modifications, ", "))
}

func (s *SchemaManagementService) buildCreateIndexSQL(tableName string, index IndexDefinition) string {
	unique := ""
	if index.Unique {
		unique = "UNIQUE "
	}

	columns := make([]string, len(index.Columns))
	for i, col := range index.Columns {
		columns[i] = pgx.Identifier{col}.Sanitize()
	}

	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)",
		unique,
		pgx.Identifier{index.Name}.Sanitize(),
		pgx.Identifier{tableName}.Sanitize(),
		strings.Join(columns, ", "))
}

func (s *SchemaManagementService) buildAddForeignKeySQL(tableName string, fk ForeignKeyDefinition) string {
	sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
		pgx.Identifier{tableName}.Sanitize(),
		pgx.Identifier{fk.Name}.Sanitize(),
		pgx.Identifier{fk.Column}.Sanitize(),
		pgx.Identifier{fk.ReferencedTable}.Sanitize(),
		pgx.Identifier{fk.ReferencedColumn}.Sanitize())

	if fk.OnDelete != "" {
		sql += fmt.Sprintf(" ON DELETE %s", strings.ToUpper(fk.OnDelete))
	}

	if fk.OnUpdate != "" {
		sql += fmt.Sprintf(" ON UPDATE %s", strings.ToUpper(fk.OnUpdate))
	}

	return sql
}
