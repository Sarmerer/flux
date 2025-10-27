package validation

import (
	"fmt"
	"strings"

	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/database"
)

type TableSchema = database.TableSchema
type ColumnDefinition = database.ColumnDefinition
type IndexDefinition = database.IndexDefinition
type ForeignKeyDefinition = database.ForeignKeyDefinition

func ValidateTableSchema(schema *TableSchema) error {
	var validationErrors []string

	if len(schema.Columns) == 0 {
		validationErrors = append(validationErrors, "at least one column is required")
	}

	columnNames := make(map[string]bool)
	for i, col := range schema.Columns {
		if err := validateColumnDefinition(col, i); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}

		if columnNames[col.Name] {
			validationErrors = append(validationErrors, fmt.Sprintf("duplicate column name at index %d", i))
		}
		columnNames[col.Name] = true
	}

	if len(schema.PrimaryKey) > 0 {
		for _, pkCol := range schema.PrimaryKey {
			if !columnNames[pkCol] {
				validationErrors = append(validationErrors, fmt.Sprintf("primary key column '%s' does not exist", pkCol))
			}
		}
	}

	for i, index := range schema.Indexes {
		if err := validateIndexDefinition(index, i); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}

		for _, colName := range index.Columns {
			if !columnNames[colName] {
				validationErrors = append(validationErrors, fmt.Sprintf("index %d: column '%s' does not exist", i, colName))
			}
		}
	}

	for i, fk := range schema.ForeignKeys {
		if err := validateForeignKeyDefinition(fk, i); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}

		if !columnNames[fk.Column] {
			validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: column '%s' does not exist", i, fk.Column))
		}
	}

	if len(validationErrors) > 0 {
		return errors.NewValidationError("Schema validation failed").
			WithField("schema").
			WithDetails(strings.Join(validationErrors, "; "))
	}

	return nil
}

func validateColumnDefinition(col ColumnDefinition, index int) error {
	var validationErrors []string

	if col.Name == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("column %d: name is required", index))
	} else if !IsValidColumnName(col.Name) {
		validationErrors = append(validationErrors, fmt.Sprintf("column %d: invalid name", index))
	}

	if col.Type == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("column %d: type is required", index))
	} else if !IsValidColumnType(col.Type) {
		validationErrors = append(validationErrors, fmt.Sprintf("column %d: invalid type", index))
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	return nil
}

func validateIndexDefinition(indexDef IndexDefinition, indexNum int) error {
	var validationErrors []string

	if indexDef.Name == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("index %d: name is required", indexNum))
	} else if !IsValidIndexName(indexDef.Name) {
		validationErrors = append(validationErrors, fmt.Sprintf("index %d: invalid name", indexNum))
	}

	if len(indexDef.Columns) == 0 {
		validationErrors = append(validationErrors, fmt.Sprintf("index %d: at least one column is required", indexNum))
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	return nil
}

func validateForeignKeyDefinition(fk ForeignKeyDefinition, index int) error {
	var validationErrors []string

	if fk.Name == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: name is required", index))
	} else if !IsValidConstraintName(fk.Name) {
		validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: invalid name", index))
	}

	if fk.Column == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: column is required", index))
	}

	if fk.ReferencedTable == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: referenced table is required", index))
	}

	if fk.ReferencedColumn == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: referenced column is required", index))
	}

	if fk.OnDelete != "" && !IsValidForeignKeyAction(fk.OnDelete) {
		validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: invalid ON DELETE action", index))
	}

	if fk.OnUpdate != "" && !IsValidForeignKeyAction(fk.OnUpdate) {
		validationErrors = append(validationErrors, fmt.Sprintf("foreign key %d: invalid ON UPDATE action", index))
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	return nil
}
