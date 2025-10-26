package validation

import (
	"fmt"

	"github.com/flow/internal/domain/entities"
)

func ValidateTableCreateRequest(req *entities.TableCreateRequest) error {
	if req.Name == "" {
		return fmt.Errorf("table name is required")
	}

	if !IsValidTableName(req.Name) {
		return fmt.Errorf("invalid table name")
	}

	if req.Schema == nil {
		return fmt.Errorf("table schema is required")
	}

	return nil
}
