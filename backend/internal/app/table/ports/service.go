package ports

import "github.com/example/flow/internal/app/table/domain"

type TableService interface {
	List() ([]domain.Table, error)
	Create(name string) (domain.Table, error)
	Update(id int64, name string) (domain.Table, error)
	Delete(id int64) error
}
