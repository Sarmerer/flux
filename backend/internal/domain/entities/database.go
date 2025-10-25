package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Database struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Host      string    `json:"host" db:"host"`
	Port      int       `json:"port" db:"port"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	Database  string    `json:"database" db:"database"`
	SSLMode   string    `json:"ssl_mode" db:"ssl_mode"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type DatabaseCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Database string `json:"database" validate:"required"`
	SSLMode  string `json:"ssl_mode"`
}

type DatabaseResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	Database  string    `json:"database"`
	SSLMode   string    `json:"ssl_mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Database) ToResponse() DatabaseResponse {
	return DatabaseResponse{
		ID:        d.ID,
		ProjectID: d.ProjectID,
		Name:      d.Name,
		Host:      d.Host,
		Port:      d.Port,
		Username:  d.Username,
		Database:  d.Database,
		SSLMode:   d.SSLMode,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (d *Database) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.SSLMode)
}
