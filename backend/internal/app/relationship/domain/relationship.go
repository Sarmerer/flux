package domain

import "time"

type Kind string

const (
	OneToOne   Kind = "o2o"
	OneToMany  Kind = "o2m"
	ManyToOne  Kind = "m2o"
	ManyToMany Kind = "m2m"
)

type Relationship struct {
	ID            int64     `json:"id"`
	SourceTableID int64     `json:"source_table_id"`
	TargetTableID int64     `json:"target_table_id"`
	Kind          Kind      `json:"kind"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
