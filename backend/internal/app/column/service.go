package column

import (
	"errors"
	"strings"

	"github.com/example/flow/internal/app/column/domain"
	"github.com/example/flow/internal/app/column/ports"
)

type Service struct {
	repo ports.ColumnRepository
}

func NewService(repo ports.ColumnRepository) *Service { return &Service{repo: repo} }

func (s *Service) List(tableID int64) ([]domain.Column, error) {
	return s.repo.List(tableID)
}

func (s *Service) Create(tableID int64, name string, typ domain.ColumnType, required bool, unique bool, defaultValue *string) (domain.Column, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Column{}, errors.New("name is required")
	}
	return s.repo.Create(tableID, name, typ, required, unique, defaultValue)
}

func (s *Service) Update(id int64, name string, required bool, unique bool, defaultValue *string) (domain.Column, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Column{}, errors.New("name is required")
	}
	return s.repo.Update(id, name, required, unique, defaultValue)
}

func (s *Service) Delete(id int64) error { return s.repo.Delete(id) }
