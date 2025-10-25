package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/flow/internal/domain/entities"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

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

func ValidateDatabaseCreateRequest(req *entities.DatabaseCreateRequest) error {
	var errors []ValidationError

	if req.Name == "" {
		errors = append(errors, ValidationError{Field: "name", Message: "name is required"})
	} else if !isValidDatabaseName(req.Name) {
		errors = append(errors, ValidationError{Field: "name", Message: "name contains invalid characters"})
	}

	if req.Host == "" {
		errors = append(errors, ValidationError{Field: "host", Message: "host is required"})
	} else if !isValidHost(req.Host) {
		errors = append(errors, ValidationError{Field: "host", Message: "invalid host format"})
	}

	if req.Port <= 0 || req.Port > 65535 {
		errors = append(errors, ValidationError{Field: "port", Message: "port must be between 1 and 65535"})
	}

	if req.Username == "" {
		errors = append(errors, ValidationError{Field: "username", Message: "username is required"})
	}

	if req.Password == "" {
		errors = append(errors, ValidationError{Field: "password", Message: "password is required"})
	}

	if req.Database == "" {
		errors = append(errors, ValidationError{Field: "database", Message: "database name is required"})
	} else if !isValidDatabaseName(req.Database) {
		errors = append(errors, ValidationError{Field: "database", Message: "database name contains invalid characters"})
	}

	if req.SSLMode != "" && !isValidSSLMode(req.SSLMode) {
		errors = append(errors, ValidationError{Field: "ssl_mode", Message: "invalid SSL mode"})
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

func ValidateTableSchema(schema *TableSchema) error {
	var errors []ValidationError

	if len(schema.Columns) == 0 {
		errors = append(errors, ValidationError{Field: "columns", Message: "at least one column is required"})
	}

	columnNames := make(map[string]bool)
	for i, col := range schema.Columns {
		if err := validateColumnDefinition(col, i); err != nil {
			if validationErr, ok := err.(ValidationErrors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d]", i), Message: err.Error()})
			}
		}

		if columnNames[col.Name] {
			errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].name", i), Message: "duplicate column name"})
		}
		columnNames[col.Name] = true
	}

	if len(schema.PrimaryKey) > 0 {
		for _, pkCol := range schema.PrimaryKey {
			if !columnNames[pkCol] {
				errors = append(errors, ValidationError{Field: "primary_key", Message: fmt.Sprintf("primary key column '%s' does not exist", pkCol)})
			}
		}
	}

	for i, index := range schema.Indexes {
		if err := validateIndexDefinition(index, i); err != nil {
			if validationErr, ok := err.(ValidationErrors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d]", i), Message: err.Error()})
			}
		}

		for _, colName := range index.Columns {
			if !columnNames[colName] {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].columns", i), Message: fmt.Sprintf("index column '%s' does not exist", colName)})
			}
		}
	}

	for i, fk := range schema.ForeignKeys {
		if err := validateForeignKeyDefinition(fk, i); err != nil {
			if validationErr, ok := err.(ValidationErrors); ok {
				errors = append(errors, validationErr.Errors...)
			} else {
				errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d]", i), Message: err.Error()})
			}
		}

		if !columnNames[fk.Column] {
			errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].column", i), Message: fmt.Sprintf("foreign key column '%s' does not exist", fk.Column)})
		}
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

func validateColumnDefinition(col ColumnDefinition, index int) error {
	var errors []ValidationError

	if col.Name == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].name", index), Message: "column name is required"})
	} else if !isValidColumnName(col.Name) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("columns[%d].name", index), Message: "invalid column name"})
	}

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

func validateIndexDefinition(indexDef IndexDefinition, indexNum int) error {
	var errors []ValidationError

	if indexDef.Name == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].name", indexNum), Message: "index name is required"})
	} else if !isValidIndexName(indexDef.Name) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].name", indexNum), Message: "invalid index name"})
	}

	if len(indexDef.Columns) == 0 {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("indexes[%d].columns", indexNum), Message: "at least one column is required"})
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

func validateForeignKeyDefinition(fk ForeignKeyDefinition, index int) error {
	var errors []ValidationError

	if fk.Name == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].name", index), Message: "foreign key name is required"})
	} else if !isValidConstraintName(fk.Name) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].name", index), Message: "invalid foreign key name"})
	}

	if fk.Column == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].column", index), Message: "column is required"})
	}

	if fk.ReferencedTable == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].referenced_table", index), Message: "referenced table is required"})
	}

	if fk.ReferencedColumn == "" {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].referenced_column", index), Message: "referenced column is required"})
	}

	if fk.OnDelete != "" && !isValidForeignKeyAction(fk.OnDelete) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].on_delete", index), Message: "invalid ON DELETE action"})
	}

	if fk.OnUpdate != "" && !isValidForeignKeyAction(fk.OnUpdate) {
		errors = append(errors, ValidationError{Field: fmt.Sprintf("foreign_keys[%d].on_update", index), Message: "invalid ON UPDATE action"})
	}

	if len(errors) > 0 {
		return ValidationErrors{Errors: errors}
	}

	return nil
}

func isValidDatabaseName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func isValidHost(host string) bool {

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
