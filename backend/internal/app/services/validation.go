package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/flow/internal/domain/entities"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e ValidationErrors) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// ValidateDatabaseCreateRequest validates a database creation request
func ValidateDatabaseCreateRequest(req *entities.DatabaseCreateRequest) error {
	var errors []ValidationError

	// Validate name
	if req.Name == "" {
		errors = append(errors, ValidationError{Field: "name", Message: "name is required"})
	} else if !isValidDatabaseName(req.Name) {
		errors = append(errors, ValidationError{Field: "name", Message: "name contains invalid characters"})
	}

	// Validate host
	if req.Host == "" {
		errors = append(errors, ValidationError{Field: "host", Message: "host is required"})
	} else if !isValidHost(req.Host) {
		errors = append(errors, ValidationError{Field: "host", Message: "invalid host format"})
	}

	// Validate port
	if req.Port <= 0 || req.Port > 65535 {
		errors = append(errors, ValidationError{Field: "port", Message: "port must be between 1 and 65535"})
	}

	// Validate username
	if req.Username == "" {
		errors = append(errors, ValidationError{Field: "username", Message: "username is required"})
	}

	// Validate password
	if req.Password == "" {
		errors = append(errors, ValidationError{Field: "password", Message: "password is required"})
	}

	// Validate database name
	if req.Database == "" {
		errors = append(errors, ValidationError{Field: "database", Message: "database name is required"})
	} else if !isValidDatabaseName(req.Database) {
		errors = append(errors, ValidationError{Field: "database", Message: "database name contains invalid characters"})
	}

	// Validate SSL mode
	if req.SSLMode != "" && !isValidSSLMode(req.SSLMode) {
		errors = append(errors, ValidationError{Field: "ssl_mode", Message: "invalid SSL mode"})
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

// ValidateTableSchema validates a table schema
func ValidateTableSchema(schema *TableSchema) error {
	var errors []ValidationError

	// Validate columns
	if len(schema.Columns) == 0 {
		errors = append(errors, ValidationError{Field: "columns", Message: "at least one column is required"})
	}

	// Validate each column
	columnNames := make(map[string]bool)
	for i, col := range schema.Columns {
		if err := validateColumnDefinition(col, i); err != nil {
			if validationErr, ok := err.(ValidationErrors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d]", i), Message: err.Error()})
			}
		}

		// Check for duplicate column names
		if columnNames[col.Name] {
			errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].name", i), Message: "duplicate column name"})
		}
		columnNames[col.Name] = true
	}

	// Validate primary key
	if len(schema.PrimaryKey) > 0 {
		for _, pkCol := range schema.PrimaryKey {
			if !columnNames[pkCol] {
				errors = append(errors, ValidationError{Field: "primary_key", Message: fmt.Sprintf("primary key column '%s' does not exist", pkCol)})
			}
		}
	}

	// Validate indexes
	for i, index := range schema.Indexes {
		if err := validateIndexDefinition(index, i); err != nil {
			if validationErr, ok := err.(ValidationErrors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d]", i), Message: err.Error()})
			}
		}

		// Check if index columns exist
		for _, colName := range index.Columns {
			if !columnNames[colName] {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].columns", i), Message: fmt.Sprintf("index column '%s' does not exist", colName)})
			}
		}
	}

	// Validate foreign keys
	for i, fk := range schema.ForeignKeys {
		if err := validateForeignKeyDefinition(fk, i); err != nil {
			if validationErr, ok := err.(ValidationErrors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d]", i), Message: err.Error()})
			}
		}

		// Check if foreign key column exists
		if !columnNames[fk.Column] {
			errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].column", i), Message: fmt.Sprintf("foreign key column '%s' does not exist", fk.Column)})
		}
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

// validateColumnDefinition validates a single column definition
func validateColumnDefinition(col ColumnDefinition, index int) error {
	var errors []ValidationError

	// Validate column name
	if col.Name == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].name", index), Message: "column name is required"})
	} else if !isValidColumnName(col.Name) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].name", index), Message: "invalid column name"})
	}

	// Validate column type
	if col.Type == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].type", index), Message: "column type is required"})
	} else if !isValidColumnType(col.Type) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].type", index), Message: "invalid column type"})
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

// validateIndexDefinition validates an index definition
func validateIndexDefinition(indexDef IndexDefinition, indexNum int) error {
	var errors []ValidationError

	// Validate index name
	if indexDef.Name == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].name", indexNum), Message: "index name is required"})
	} else if !isValidIndexName(indexDef.Name) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].name", indexNum), Message: "invalid index name"})
	}

	// Validate columns
	if len(indexDef.Columns) == 0 {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].columns", indexNum), Message: "at least one column is required"})
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

// validateForeignKeyDefinition validates a foreign key definition
func validateForeignKeyDefinition(fk ForeignKeyDefinition, index int) error {
	var errors []ValidationError

	// Validate foreign key name
	if fk.Name == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].name", index), Message: "foreign key name is required"})
	} else if !isValidConstraintName(fk.Name) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].name", index), Message: "invalid foreign key name"})
	}

	// Validate column
	if fk.Column == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].column", index), Message: "column is required"})
	}

	// Validate referenced table
	if fk.ReferencedTable == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].referenced_table", index), Message: "referenced table is required"})
	}

	// Validate referenced column
	if fk.ReferencedColumn == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].referenced_column", index), Message: "referenced column is required"})
	}

	// Validate ON DELETE action
	if fk.OnDelete != "" && !isValidForeignKeyAction(fk.OnDelete) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].on_delete", index), Message: "invalid ON DELETE action"})
	}

	// Validate ON UPDATE action
	if fk.OnUpdate != "" && !isValidForeignKeyAction(fk.OnUpdate) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].on_update", index), Message: "invalid ON UPDATE action"})
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

// Validation helper functions

func isValidDatabaseName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func isValidHost(host string) bool {
	// Basic host validation - can be IP or hostname
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9.-]+$`, host)
	return matched && len(host) <= 255
}

func isValidSSLMode(mode string) bool {
	validModes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	for _, validMode := range validModes {
		if mode == validMode {
			return true
		}
	}
	return false
}

func isValidColumnName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func isValidColumnType(colType string) bool {
	validTypes := []string{
		"SMALLINT", "INTEGER", "BIGINT", "DECIMAL", "NUMERIC", "REAL", "DOUBLE PRECISION",
		"SMALLSERIAL", "SERIAL", "BIGSERIAL", "MONEY",
		"CHARACTER VARYING", "VARCHAR", "CHARACTER", "CHAR", "TEXT",
		"BYTEA", "TIMESTAMP", "TIMESTAMPTZ", "DATE", "TIME", "TIMETZ", "INTERVAL",
		"BOOLEAN", "POINT", "LINE", "LSEG", "BOX", "PATH", "POLYGON", "CIRCLE",
		"CIDR", "INET", "MACADDR", "UUID", "XML", "JSON", "JSONB", "ARRAY",
	}

	colType = strings.ToUpper(colType)
	for _, validType := range validTypes {
		if strings.HasPrefix(colType, validType) {
			return true
		}
	}
	return false
}

func isValidIndexName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func isValidConstraintName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func isValidForeignKeyAction(action string) bool {
	validActions := []string{"NO ACTION", "RESTRICT", "CASCADE", "SET NULL", "SET DEFAULT"}
	action = strings.ToUpper(action)
	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}

func isValidTableName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}
