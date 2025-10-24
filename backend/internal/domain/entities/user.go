package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Role represents a user role
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleManager   Role = "manager"
	RoleDeveloper Role = "developer"
	RoleViewer    Role = "viewer"
)

// Permission represents a user permission
type Permission string

const (
	PermProjectsCreate Permission = "projects.create"
	PermProjectsEdit   Permission = "projects.edit"
	PermProjectsDelete Permission = "projects.delete"
	PermProjectsView   Permission = "projects.view"
	PermWorkflowsCreate Permission = "workflows.create"
	PermWorkflowsEdit   Permission = "workflows.edit"
	PermWorkflowsDelete Permission = "workflows.delete"
	PermTablesCreate    Permission = "tables.create"
	PermTablesEdit      Permission = "tables.edit"
	PermTablesDelete    Permission = "tables.delete"
	PermSettingsManage  Permission = "settings.manage"
	PermUsersManage     Permission = "users.manage"
)

// Permissions is a slice of Permission that can be stored in JSONB
type Permissions []Permission

// Value implements the driver.Valuer interface for database storage
func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for database retrieval
func (p *Permissions) Scan(value interface{}) error {
	if value == nil {
		*p = []Permission{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, p)
}

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserCreateRequest represents the data needed to create a new user
type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

// UserLoginRequest represents the data needed for user login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts a User entity to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

