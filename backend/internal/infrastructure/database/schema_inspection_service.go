package database

import (
	"context"

	"github.com/flow/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchemaInspectionService struct{}

func NewSchemaInspectionService() *SchemaInspectionService {
	return &SchemaInspectionService{}
}

func (s *SchemaInspectionService) GetTableSchema(
	ctx context.Context,
	db *pgxpool.Pool,
	tableName string,
) (*entities.TableSchemaInfo, error) {
	columns, err := s.getColumns(ctx, db, tableName)
	if err != nil {
		return nil, err
	}

	primaryKeys, err := s.getPrimaryKeys(ctx, db, tableName)
	if err != nil {
		return nil, err
	}

	foreignKeys, err := s.getForeignKeys(ctx, db, tableName)
	if err != nil {
		return nil, err
	}

	indexes, err := s.getIndexes(ctx, db, tableName)
	if err != nil {
		return nil, err
	}

	return &entities.TableSchemaInfo{
		Columns:     columns,
		PrimaryKeys: primaryKeys,
		ForeignKeys: foreignKeys,
		Indexes:     indexes,
	}, nil
}

func (s *SchemaInspectionService) getColumns(
	ctx context.Context,
	db *pgxpool.Pool,
	tableName string,
) ([]entities.ColumnInfo, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable = 'YES' as is_nullable,
			column_default,
			character_maximum_length,
			numeric_precision,
			numeric_scale,
			ordinal_position
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.Query(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []entities.ColumnInfo
	for rows.Next() {
		var col entities.ColumnInfo
		err := rows.Scan(
			&col.Name,
			&col.DataType,
			&col.IsNullable,
			&col.DefaultValue,
			&col.MaxLength,
			&col.NumericPrecision,
			&col.NumericScale,
			&col.OrdinalPosition,
		)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, rows.Err()
}

func (s *SchemaInspectionService) getPrimaryKeys(
	ctx context.Context,
	db *pgxpool.Pool,
	tableName string,
) ([]string, error) {
	query := `
		SELECT kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name = $1
			AND tc.constraint_type = 'PRIMARY KEY'
		ORDER BY kcu.ordinal_position
	`

	rows, err := db.Query(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pks []string
	for rows.Next() {
		var colName string
		if err := rows.Scan(&colName); err != nil {
			return nil, err
		}
		pks = append(pks, colName)
	}

	return pks, rows.Err()
}

func (s *SchemaInspectionService) getForeignKeys(
	ctx context.Context,
	db *pgxpool.Pool,
	tableName string,
) ([]entities.ForeignKeyInfo, error) {
	query := `
		SELECT
			tc.constraint_name,
			kcu.column_name,
			ccu.table_name AS referenced_table,
			ccu.column_name AS referenced_column,
			rc.update_rule,
			rc.delete_rule
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage ccu
			ON tc.constraint_name = ccu.constraint_name
		JOIN information_schema.referential_constraints rc
			ON tc.constraint_name = rc.constraint_name
		WHERE tc.table_name = $1
			AND tc.constraint_type = 'FOREIGN KEY'
		ORDER BY tc.constraint_name, kcu.ordinal_position
	`

	rows, err := db.Query(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fkMap := make(map[string]*entities.ForeignKeyInfo)

	for rows.Next() {
		var (
			constraintName, colName, refTable, refCol, updateRule, deleteRule string
		)
		err := rows.Scan(&constraintName, &colName, &refTable, &refCol, &updateRule, &deleteRule)
		if err != nil {
			return nil, err
		}

		if fk, exists := fkMap[constraintName]; exists {
			fk.ColumnNames = append(fk.ColumnNames, colName)
			fk.ReferencedColumns = append(fk.ReferencedColumns, refCol)
		} else {
			fkMap[constraintName] = &entities.ForeignKeyInfo{
				ConstraintName:    constraintName,
				ColumnNames:       []string{colName},
				ReferencedTable:   refTable,
				ReferencedColumns: []string{refCol},
				UpdateRule:        updateRule,
				DeleteRule:        deleteRule,
			}
		}
	}

	var fks []entities.ForeignKeyInfo
	for _, fk := range fkMap {
		fks = append(fks, *fk)
	}

	return fks, rows.Err()
}

func (s *SchemaInspectionService) getIndexes(
	ctx context.Context,
	db *pgxpool.Pool,
	tableName string,
) ([]entities.IndexInfo, error) {
	query := `
		SELECT
			i.relname AS index_name,
			a.attname AS column_name,
			ix.indisunique AS is_unique,
			ix.indisprimary AS is_primary,
			am.amname AS index_type
		FROM pg_class t
		JOIN pg_index ix ON t.oid = ix.indrelid
		JOIN pg_class i ON i.oid = ix.indexrelid
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		JOIN pg_am am ON i.relam = am.oid
		WHERE t.relname = $1
		ORDER BY i.relname, a.attnum
	`

	rows, err := db.Query(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexMap := make(map[string]*entities.IndexInfo)

	for rows.Next() {
		var (
			indexName, colName, indexType string
			isUnique, isPrimary           bool
		)
		err := rows.Scan(&indexName, &colName, &isUnique, &isPrimary, &indexType)
		if err != nil {
			return nil, err
		}

		if idx, exists := indexMap[indexName]; exists {
			idx.ColumnNames = append(idx.ColumnNames, colName)
		} else {
			indexMap[indexName] = &entities.IndexInfo{
				Name:        indexName,
				ColumnNames: []string{colName},
				IsUnique:    isUnique,
				IsPrimary:   isPrimary,
				IndexType:   indexType,
			}
		}
	}

	var indexes []entities.IndexInfo
	for _, idx := range indexMap {
		indexes = append(indexes, *idx)
	}

	return indexes, rows.Err()
}
