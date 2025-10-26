package validation

import "github.com/flow/internal/domain/entities"

func ValidateDatabaseCreateRequest(req *entities.DatabaseCreateRequest) error {
	var errors []Error

	if req.Name == "" {
		errors = append(errors, Error{Field: "name", Message: "name is required"})
	} else if !IsValidDatabaseName(req.Name) {
		errors = append(errors, Error{Field: "name", Message: "name contains invalid characters"})
	}

	if req.Host == "" {
		errors = append(errors, Error{Field: "host", Message: "host is required"})
	} else if !IsValidHost(req.Host) {
		errors = append(errors, Error{Field: "host", Message: "invalid host format"})
	}

	if req.Port <= 0 || req.Port > 65535 {
		errors = append(errors, Error{Field: "port", Message: "port must be between 1 and 65535"})
	}

	if req.Username == "" {
		errors = append(errors, Error{Field: "username", Message: "username is required"})
	}

	if req.Password == "" {
		errors = append(errors, Error{Field: "password", Message: "password is required"})
	}

	if req.Database == "" {
		errors = append(errors, Error{Field: "database", Message: "database name is required"})
	} else if !IsValidDatabaseName(req.Database) {
		errors = append(errors, Error{Field: "database", Message: "database name contains invalid characters"})
	}

	if req.SSLMode != "" && !IsValidSSLMode(req.SSLMode) {
		errors = append(errors, Error{Field: "ssl_mode", Message: "invalid SSL mode"})
	}

	if len(errors) > 0 {
		return Errors{Errors: errors}
	}

	return nil
}
