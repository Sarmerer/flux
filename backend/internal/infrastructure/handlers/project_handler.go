package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectService services.ProjectServiceInterface
}

func NewProjectHandler(projectService services.ProjectServiceInterface) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req entities.ProjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.NewValidationError("Invalid request body").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		errors.WriteError(w, errors.NewValidationError("Project name is required").WithField("name"))
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError())
		return
	}

	project, err := h.projectService.CreateProject(r.Context(), &req, userID)
	if err != nil {

		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewInternalError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid project ID").WithField("id").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	project, err := h.projectService.GetProjectByID(r.Context(), projectID)
	if err != nil {

		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewNotFoundError("Project"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError())
		return
	}

	projects, err := h.projectService.GetProjectsByOwnerID(r.Context(), userID)
	if err != nil {

		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewInternalError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid project ID").WithField("id").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	var req entities.ProjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.NewValidationError("Invalid request body").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		errors.WriteError(w, errors.NewValidationError("Project name is required").WithField("name"))
		return
	}

	project, err := h.projectService.UpdateProject(r.Context(), projectID, &req)
	if err != nil {

		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewInternalError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid project ID").WithField("id").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	if err := h.projectService.DeleteProject(r.Context(), projectID); err != nil {

		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewInternalError(err))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
