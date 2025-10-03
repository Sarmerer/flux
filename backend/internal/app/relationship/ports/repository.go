package ports

import "github.com/example/flow/internal/app/relationship/domain"

type Repository interface {
	List(tableID int64) ([]domain.Relationship, error)
	Create(sourceID, targetID int64, kind domain.Kind) (domain.Relationship, error)
	Delete(id int64) error
}
