package validation

import (
	"fmt"

	"github.com/flow/internal/infrastructure/database"
)

type TableSchema = database.TableSchema
type ColumnDefinition = database.ColumnDefinition
type IndexDefinition = database.IndexDefinition
type ForeignKeyDefinition = database.ForeignKeyDefinition

func ValidateTableSchema(schema *TableSchema) error {
	var errors []Error

	if len(schema.Columns) == 0 {
		errors = append(errors, Error{Field: "columns", Message: "at least one column is required"})
	}

	columnNames := make(map[string]bool)
	for i, col := range schema.Columns {
		if err := validateColumnDefinition(col, i); err != nil {
			if validationErr, ok := err.(Errors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, Error{Field: fmt.Sprintf("columns[%d]", i), Message: err.Error()})
			}
		}

		if columnNames[col.Name] {
			errors = append(errors, Error{Field: fmt.Sprintf("columns[%d].name", i), Message: "duplicate column name"})
		}
		columnNames[col.Name] = true
	}

	if len(schema.PrimaryKey) > 0 {
		for _, pkCol := range schema.PrimaryKey {
			if !columnNames[pkCol] {
				errors = append(errors, Error{Field: "primary_key", Message: fmt.Sprintf("primary key column '%s' does not exist", pkCol)})
			}
		}
	}

	for i, index := range schema.Indexes {
		if err := validateIndexDefinition(index, i); err != nil {
			if validationErr, ok := err.(Errors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, Error{Field: fmt.Sprintf("indexes[%d]", i), Message: err.Error()})
			}
		}

		for _, colName := range index.Columns {
			if !columnNames[colName] {
				errors = append(errors, Error{Field: fmt.Sprintf("indexes[%d].columns", i), Message: fmt.Sprintf("index column '%s' does not exist", colName)})
			}
		}
	}

	for i, fk := range schema.ForeignKeys {
		if err := validateForeignKeyDefinition(fk, i); err != nil {
			if validationErr, ok := err.(Errors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d]", i), Message: err.Error()})
			}
		}

		if !columnNames[fk.Column] {
			errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].column", i), Message: fmt.Sprintf("foreign key column '%s' does not exist", fk.Column)})
		}
	}

	if len(errors) > 0 {
		return Errors{Errors: errors}
	}

	return nil
}

func validateColumnDefinition(col ColumnDefinition, index int) error {
	var errors []Error

	if col.Name == "" {
		errors = append(errors, Error{Field: fmt.Sprintf("columns[%d].name", index), Message: "column name is required"})
	} else if !IsValidColumnName(col.Name) {
		errors = append(errors, Error{Field: fmt.Sprintf("columns[%d].name", index), Message: "invalid column name"})
	}

	if col.Type == "" {
		errors = append(errors, Error{Field: fmt.Sprintf("columns[%d].type", index), Message: "column type is required"})
	} else if !IsValidColumnType(col.Type) {
		errors = append(errors, Error{Field: fmt.Sprintf("columns[%d].type", index), Message: "invalid column type"})
	}

	if len(errors) > 0 {
		return Errors{Errors: errors}
	}

	return nil
}

func validateIndexDefinition(indexDef IndexDefinition, indexNum int) error {
	var errors []Error

	if indexDef.Name == "" {
		errors = append(errors, Error{Field: fmt.Sprintf("indexes[%d].name", indexNum), Message: "index name is required"})
	} else if !IsValidIndexName(indexDef.Name) {
		errors = append(errors, Error{Field: fmt.Sprintf("indexes[%d].name", indexNum), Message: "invalid index name"})
	}

	if len(indexDef.Columns) == 0 {
		errors = append(errors, Error{Field: fmt.Sprintf("indexes[%d].columns", indexNum), Message: "at least one column is required"})
	}

	if len(errors) > 0 {
		return Errors{Errors: errors}
	}

	return nil
}

func validateForeignKeyDefinition(fk ForeignKeyDefinition, index int) error {
	var errors []Error

	if fk.Name == "" {
		errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].name", index), Message: "foreign key name is required"})
	} else if !IsValidConstraintName(fk.Name) {
		errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].name", index), Message: "invalid foreign key name"})
	}

	if fk.Column == "" {
		errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].column", index), Message: "column is required"})
	}

	if fk.ReferencedTable == "" {
		errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].referenced_table", index), Message: "referenced table is required"})
	}

	if fk.ReferencedColumn == "" {
		errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].referenced_column", index), Message: "referenced column is required"})
	}

	if fk.OnDelete != "" && !IsValidForeignKeyAction(fk.OnDelete) {
		errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].on_delete", index), Message: "invalid ON DELETE action"})
	}

	if fk.OnUpdate != "" && !IsValidForeignKeyAction(fk.OnUpdate) {
		errors = append(errors, Error{Field: fmt.Sprintf("foreign_keys[%d].on_update", index), Message: "invalid ON UPDATE action"})
	}

	if len(errors) > 0 {
		return Errors{Errors: errors}
	}

	return nil
}
