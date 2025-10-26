package validation

import (
	"regexp"
	"strings"
)

func IsValidDatabaseName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func IsValidHost(host string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9.-]+$`, host)
	return matched && len(host) <= 255
}

func IsValidSSLMode(mode string) bool {
	validModes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	for _, validMode := range validModes {
		if mode == validMode {
			return true
		}
	}
	return false
}

func IsValidColumnName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func IsValidTableName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func IsValidIndexName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func IsValidConstraintName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched && len(name) <= 63
}

func IsValidForeignKeyAction(action string) bool {
	validActions := []string{"NO ACTION", "RESTRICT", "CASCADE", "SET NULL", "SET DEFAULT"}
	actionUpper := strings.ToUpper(action)
	for _, validAction := range validActions {
		if actionUpper == validAction {
			return true
		}
	}
	return false
}

func IsValidColumnType(colType string) bool {
	validTypes := []string{
		"SMALLINT", "INTEGER", "BIGINT", "DECIMAL", "NUMERIC", "REAL", "DOUBLE PRECISION",
		"SMALLSERIAL", "SERIAL", "BIGSERIAL", "MONEY",
		"CHARACTER VARYING", "VARCHAR", "CHARACTER", "CHAR", "TEXT",
		"BYTEA", "TIMESTAMP", "TIMESTAMPTZ", "DATE", "TIME", "TIMETZ", "INTERVAL",
		"BOOLEAN", "POINT", "LINE", "LSEG", "BOX", "PATH", "POLYGON", "CIRCLE",
		"CIDR", "INET", "MACADDR", "UUID", "XML", "JSON", "JSONB", "ARRAY",
	}

	colTypeUpper := strings.ToUpper(colType)
	for _, validType := range validTypes {
		if strings.HasPrefix(colTypeUpper, validType) {
			return true
		}
	}
	return false
}
