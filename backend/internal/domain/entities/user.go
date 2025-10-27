package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleManager   Role = "manager"
	RoleDeveloper Role = "developer"
	RoleViewer    Role = "viewer"
)

type Permission string

const (
	PermProjectCreate Permission = "project.create"
	PermProjectEdit   Permission = "project.edit"
	PermProjectDelete Permission = "project.delete"
	PermProjectView   Permission = "project.view"

	PermWorkflowCreate Permission = "workflow.create"
	PermWorkflowEdit   Permission = "workflow.edit"
	PermWorkflowDelete Permission = "workflow.delete"

	PermTableCreate Permission = "table.create"
	PermTableEdit   Permission = "table.edit"
	PermTableDelete Permission = "table.delete"

	PermSettingManage Permission = "setting.manage"
	PermUserManage    Permission = "user.manage"
)

type Permissions []Permission

func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

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

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
