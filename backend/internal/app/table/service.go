package table

import (
	"errors"
	"strings"

	"github.com/example/flow/internal/app/table/domain"
	"github.com/example/flow/internal/app/table/ports"
)

type Service struct {
	repo ports.TableRepository
}

func NewService(repo ports.TableRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]domain.Table, error) {
	return s.repo.ListTables()
}

func (s *Service) Create(name string) (domain.Table, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Table{}, errors.New("name is required")
	}
	return s.repo.CreateTable(name)
}

func (s *Service) Update(id int64, name string) (domain.Table, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Table{}, errors.New("name is required")
	}
	return s.repo.UpdateTable(id, name)
}

func (s *Service) Delete(id int64) error {
	return s.repo.DeleteTable(id)
}
