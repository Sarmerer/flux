package validation

import (
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
)

func ValidateTableCreateRequest(req *entities.TableCreateRequest) error {
	if req.Name == "" {
		return errors.NewValidationError("Table name is required").WithField("name")
	}

	if !IsValidTableName(req.Name) {
		return errors.NewValidationError("Invalid table name").WithField("name")
	}

	if req.Schema == nil {
		return errors.NewValidationError("Table schema is required").WithField("schema")
	}

	return nil
}
