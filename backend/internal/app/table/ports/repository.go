package ports

import "github.com/example/flow/internal/app/table/domain"

type TableRepository interface {
	ListTables() ([]domain.Table, error)
	CreateTable(name string) (domain.Table, error)
	UpdateTable(id int64, name string) (domain.Table, error)
	DeleteTable(id int64) error
}
