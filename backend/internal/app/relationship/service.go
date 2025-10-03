package relationship

import (
	"errors"

	"github.com/example/flow/internal/app/relationship/domain"
	"github.com/example/flow/internal/app/relationship/ports"
)

type Service struct{ repo ports.Repository }

func NewService(r ports.Repository) *Service { return &Service{repo: r} }

func (s *Service) List(tableID int64) ([]domain.Relationship, error) { return s.repo.List(tableID) }

func (s *Service) Create(sourceID, targetID int64, kind domain.Kind) (domain.Relationship, error) {
	if sourceID == 0 || targetID == 0 {
		return domain.Relationship{}, errors.New("table ids required")
	}
	switch kind {
	case domain.OneToOne, domain.OneToMany, domain.ManyToOne, domain.ManyToMany:
	default:
		return domain.Relationship{}, errors.New("invalid kind")
	}
	return s.repo.Create(sourceID, targetID, kind)
}

func (s *Service) Delete(id int64) error { return s.repo.Delete(id) }
