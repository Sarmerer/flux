package handlers

import (
	"net/http"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/middleware"
	"github.com/google/uuid"
)

type ProjectMemberHandler struct {
	memberRepo  repositories.ProjectMemberRepository
	userRepo    repositories.UserRepository
	projectRepo repositories.ProjectRepository
}

func NewProjectMemberHandler(
	memberRepo repositories.ProjectMemberRepository,
	userRepo repositories.UserRepository,
	projectRepo repositories.ProjectRepository,
) *ProjectMemberHandler {
	return &ProjectMemberHandler{
		memberRepo:  memberRepo,
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

type AddMemberRequest struct {
	Email       string               `json:"email"`
	Role        entities.Role        `json:"role"`
	Permissions entities.Permissions `json:"permissions"`
}

type UpdateMemberRequest struct {
	Role        entities.Role        `json:"role"`
	Permissions entities.Permissions `json:"permissions"`
}

type ProjectMemberWithUser struct {
	ID          uuid.UUID             `json:"id"`
	ProjectID   uuid.UUID             `json:"project_id"`
	User        entities.UserResponse `json:"user"`
	Role        entities.Role         `json:"role"`
	Permissions entities.Permissions  `json:"permissions"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

func (h *ProjectMemberHandler) GetProjectMembers(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	members, err := h.memberRepo.GetByProjectID(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewDatabaseError("Failed to get project members", err))
		return
	}

	response := make([]ProjectMemberWithUser, 0)
	for _, member := range members {
		user, err := h.userRepo.GetByID(r.Context(), member.UserID)
		if err != nil {
			continue
		}

		permissions := member.Permissions
		if permissions == nil {
			permissions = entities.Permissions{}
		}

		response = append(response, ProjectMemberWithUser{
			ID:          member.ID,
			ProjectID:   member.ProjectID,
			User:        user.ToResponse(),
			Role:        member.Role,
			Permissions: permissions,
			CreatedAt:   member.CreatedAt,
			UpdatedAt:   member.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *ProjectMemberHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	var req AddMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if _, err := h.projectRepo.GetByID(r.Context(), projectID); err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project not found"))
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("User not found"))
		return
	}

	existingMember, _ := h.memberRepo.GetByProjectAndUser(r.Context(), projectID, user.ID)
	if existingMember != nil {
		errors.WriteError(w, errors.NewAlreadyExistsError("User is already a member of this project already exists"))
		return
	}

	permissions := req.Permissions
	if permissions == nil {
		permissions = entities.Permissions{}
	}

	member := &entities.ProjectMember{
		ID:          uuid.New(),
		ProjectID:   projectID,
		UserID:      user.ID,
		Role:        req.Role,
		Permissions: permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.memberRepo.Create(r.Context(), member); err != nil {
		errors.WriteError(w, errors.NewDatabaseError("Failed to add project member", err))
		return
	}

	responsePermissions := member.Permissions
	if responsePermissions == nil {
		responsePermissions = entities.Permissions{}
	}

	response := ProjectMemberWithUser{
		ID:          member.ID,
		ProjectID:   member.ProjectID,
		User:        user.ToResponse(),
		Role:        member.Role,
		Permissions: responsePermissions,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *ProjectMemberHandler) UpdateProjectMember(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	memberID, err := parseUUIDParam(r, "memberId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid member ID"))
		return
	}

	var req UpdateMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	member, err := h.memberRepo.GetByID(r.Context(), memberID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project member not found"))
		return
	}

	if member.ProjectID != projectID {
		errors.WriteError(w, errors.NewValidationError("Member does not belong to this project"))
		return
	}

	member.Role = req.Role

	permissions := req.Permissions
	if permissions == nil {
		permissions = entities.Permissions{}
	}
	member.Permissions = permissions
	member.UpdatedAt = time.Now()

	if err := h.memberRepo.Update(r.Context(), member); err != nil {
		errors.WriteError(w, errors.NewDatabaseError("Failed to update project member", err))
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), member.UserID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("User not found"))
		return
	}

	responsePermissions := member.Permissions
	if responsePermissions == nil {
		responsePermissions = entities.Permissions{}
	}

	response := ProjectMemberWithUser{
		ID:          member.ID,
		ProjectID:   member.ProjectID,
		User:        user.ToResponse(),
		Role:        member.Role,
		Permissions: responsePermissions,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *ProjectMemberHandler) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	memberID, err := parseUUIDParam(r, "memberId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid member ID"))
		return
	}

	member, err := h.memberRepo.GetByID(r.Context(), memberID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project member not found"))
		return
	}

	if member.ProjectID != projectID {
		errors.WriteError(w, errors.NewValidationError("Member does not belong to this project"))
		return
	}

	currentUserID, _ := middleware.GetUserIDFromContext(r.Context())
	if member.UserID == currentUserID {
		errors.WriteError(w, errors.NewValidationError("Cannot remove yourself from the project"))
		return
	}

	if err := h.memberRepo.Delete(r.Context(), memberID); err != nil {
		errors.WriteError(w, errors.NewDatabaseError("Failed to remove project member", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectMemberHandler) GetMyProjectRole(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	member, err := h.memberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project member not found"))
		return
	}

	response := member.ToResponse()
	writeJSON(w, http.StatusOK, response)
}
