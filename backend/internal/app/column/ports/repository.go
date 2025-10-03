package ports

import "github.com/example/flow/internal/app/column/domain"

type ColumnRepository interface {
	List(tableID int64) ([]domain.Column, error)
	Create(tableID int64, name string, typ domain.ColumnType, required bool, unique bool, defaultValue *string) (domain.Column, error)
	Update(id int64, name string, required bool, unique bool, defaultValue *string) (domain.Column, error)
	Delete(id int64) error
}
