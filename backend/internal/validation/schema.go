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

	if col.IsIdentity {
		if !isNumericType(col.Type) {
			validationErrors = append(validationErrors, fmt.Sprintf("column %d: identity columns must be numeric (SMALLINT, INTEGER, or BIGINT)", index))
		}
		if col.DefaultValue != "" {
			validationErrors = append(validationErrors, fmt.Sprintf("column %d: identity columns cannot have default values", index))
		}
	}

	if col.DefaultValue != "" {
		if err := ValidateDefaultExpression(col.DefaultValue, col.Type); err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("column %d: %s", index, err.Error()))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	return nil
}

func isNumericType(columnType string) bool {
	upperType := strings.ToUpper(columnType)
	numericTypes := []string{"SMALLINT", "INT2", "INTEGER", "INT", "INT4", "BIGINT", "INT8"}
	for _, t := range numericTypes {
		if upperType == t {
			return true
		}
	}
	return false
}

func ValidateDefaultExpression(defaultValue, columnType string) error {
	trimmed := strings.TrimSpace(defaultValue)
	if trimmed == "" {
		return nil
	}

	upperDefault := strings.ToUpper(trimmed)
	upperType := strings.ToUpper(columnType)

	validFunctions := []string{
		"NOW()", "CURRENT_TIMESTAMP", "CURRENT_DATE", "CURRENT_TIME",
		"GEN_RANDOM_UUID()", "UUID_GENERATE_V4()",
		"CURRENT_USER", "SESSION_USER", "CURRENT_SCHEMA",
		"TRUE", "FALSE", "NULL",
	}

	for _, fn := range validFunctions {
		if upperDefault == fn {
			return validateDefaultTypeCompatibility(fn, upperType)
		}
	}

	if strings.HasPrefix(trimmed, "'") && strings.HasSuffix(trimmed, "'") {
		return nil
	}

	if trimmed == "0" || trimmed == "1" {
		return nil
	}

	if _, err := fmt.Sscanf(trimmed, "%f", new(float64)); err == nil {
		return nil
	}

	if strings.Contains(upperDefault, "NEXTVAL") ||
		strings.Contains(upperDefault, "CURRVAL") ||
		strings.Contains(upperDefault, "RANDOM") ||
		strings.Contains(upperDefault, "EXTRACT") ||
		strings.Contains(upperDefault, "LOCALTIMESTAMP") {
		return nil
	}

	return fmt.Errorf("invalid default expression '%s'", defaultValue)
}

func validateDefaultTypeCompatibility(function, columnType string) error {
	timestampFunctions := []string{"NOW()", "CURRENT_TIMESTAMP", "LOCALTIMESTAMP"}
	dateFunctions := []string{"CURRENT_DATE"}
	timeFunctions := []string{"CURRENT_TIME"}
	uuidFunctions := []string{"GEN_RANDOM_UUID()", "UUID_GENERATE_V4()"}
	boolFunctions := []string{"TRUE", "FALSE"}

	for _, fn := range timestampFunctions {
		if function == fn && !strings.Contains(columnType, "TIMESTAMP") {
			return fmt.Errorf("%s is only valid for TIMESTAMP columns", function)
		}
	}

	for _, fn := range dateFunctions {
		if function == fn && columnType != "DATE" {
			return fmt.Errorf("%s is only valid for DATE columns", function)
		}
	}

	for _, fn := range timeFunctions {
		if function == fn && !strings.Contains(columnType, "TIME") {
			return fmt.Errorf("%s is only valid for TIME columns", function)
		}
	}

	for _, fn := range uuidFunctions {
		if function == fn && columnType != "UUID" {
			return fmt.Errorf("%s is only valid for UUID columns", function)
		}
	}

	for _, fn := range boolFunctions {
		if function == fn && columnType != "BOOLEAN" && columnType != "BOOL" {
			return fmt.Errorf("%s is only valid for BOOLEAN columns", function)
		}
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
