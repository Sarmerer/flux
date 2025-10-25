package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	ProjectID   uuid.UUID   `json:"project_id" db:"project_id"`
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	Role        Role        `json:"role" db:"role"`
	Permissions Permissions `json:"permissions" db:"permissions"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type ProjectMemberResponse struct {
	ID          uuid.UUID   `json:"id"`
	ProjectID   uuid.UUID   `json:"project_id"`
	UserID      uuid.UUID   `json:"user_id"`
	Role        Role        `json:"role"`
	Permissions Permissions `json:"permissions"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (pm *ProjectMember) ToResponse() ProjectMemberResponse {
	permissions := pm.Permissions
	if permissions == nil {
		permissions = Permissions{}
	}
	return ProjectMemberResponse{
		ID:          pm.ID,
		ProjectID:   pm.ProjectID,
		UserID:      pm.UserID,
		Role:        pm.Role,
		Permissions: permissions,
		CreatedAt:   pm.CreatedAt,
		UpdatedAt:   pm.UpdatedAt,
	}
}

func (pm *ProjectMember) HasPermission(perm Permission) bool {

	if pm.Role == RoleAdmin {
		return true
	}

	for _, p := range pm.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

func (pm *ProjectMember) HasAnyPermission(perms ...Permission) bool {
	for _, perm := range perms {
		if pm.HasPermission(perm) {
			return true
		}
	}
	return false
}

func (pm *ProjectMember) HasAllPermissions(perms ...Permission) bool {
	for _, perm := range perms {
		if !pm.HasPermission(perm) {
			return false
		}
	}
	return true
}

func (pm *ProjectMember) HasRole(roles ...Role) bool {
	for _, role := range roles {
		if pm.Role == role {
			return true
		}
	}
	return false
}

func (pm *ProjectMember) IsAdmin() bool {
	return pm.Role == RoleAdmin
}
