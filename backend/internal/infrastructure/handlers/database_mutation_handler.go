package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// DatabaseMutationHandler handles database mutation HTTP requests
type DatabaseMutationHandler struct {
	dbMutationService *services.DatabaseMutationService
}

// NewDatabaseMutationHandler creates a new DatabaseMutationHandler
func NewDatabaseMutationHandler(dbMutationService *services.DatabaseMutationService) *DatabaseMutationHandler {
	return &DatabaseMutationHandler{
		dbMutationService: dbMutationService,
	}
}

// CreateProjectDatabase handles project database creation
func (h *DatabaseMutationHandler) CreateProjectDatabase(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid project ID format").WithField("projectId")
		errors.WriteError(w, apiErr)
		return
	}

	var req entities.DatabaseCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.NewValidationError("Invalid request body").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	database, err := h.dbMutationService.CreateProjectDatabase(r.Context(), projectID, &req)
	if err != nil {
		apiErr := errors.NewDatabaseError(err)
		errors.WriteError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(database)
}

// UpdateProjectDatabase handles project database updates
func (h *DatabaseMutationHandler) UpdateProjectDatabase(w http.ResponseWriter, r *http.Request) {
	databaseIDStr := chi.URLParam(r, "id")
	databaseID, err := uuid.Parse(databaseIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid database ID format").WithField("id")
		errors.WriteError(w, apiErr)
		return
	}

	var req entities.DatabaseCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.NewValidationError("Invalid request body").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	database, err := h.dbMutationService.UpdateProjectDatabase(r.Context(), databaseID, &req)
	if err != nil {
		apiErr := errors.NewDatabaseError(err)
		errors.WriteError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(database)
}

// DeleteProjectDatabase handles project database deletion
func (h *DatabaseMutationHandler) DeleteProjectDatabase(w http.ResponseWriter, r *http.Request) {
	databaseIDStr := chi.URLParam(r, "id")
	databaseID, err := uuid.Parse(databaseIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid database ID format").WithField("id")
		errors.WriteError(w, apiErr)
		return
	}

	if err := h.dbMutationService.DeleteProjectDatabase(r.Context(), databaseID); err != nil {
		apiErr := errors.NewDatabaseError(err)
		errors.WriteError(w, apiErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TestDatabaseConnection handles database connection testing
func (h *DatabaseMutationHandler) TestDatabaseConnection(w http.ResponseWriter, r *http.Request) {
	databaseIDStr := chi.URLParam(r, "id")
	databaseID, err := uuid.Parse(databaseIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid database ID format").WithField("id")
		errors.WriteError(w, apiErr)
		return
	}

	if err := h.dbMutationService.TestDatabaseConnection(r.Context(), databaseID); err != nil {
		apiErr := errors.NewDatabaseError(err)
		errors.WriteError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "connected"})
}

// GetProjectDatabases handles getting all databases for a project
func (h *DatabaseMutationHandler) GetProjectDatabases(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid project ID format").WithField("projectId")
		errors.WriteError(w, apiErr)
		return
	}

	databases, err := h.dbMutationService.GetProjectDatabases(r.Context(), projectID)
	if err != nil {
		apiErr := errors.NewDatabaseError(err)
		errors.WriteError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(databases)
}
