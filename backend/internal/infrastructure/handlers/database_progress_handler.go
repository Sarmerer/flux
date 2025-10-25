package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
	authMiddleware "github.com/flow/internal/infrastructure/middleware"
	"github.com/flow/internal/infrastructure/progress"
	"github.com/flow/internal/infrastructure/realtime"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DatabaseMutationProgressHandler struct {
	dbService       *services.DatabaseService
	progressTracker *progress.Tracker
	hub             *realtime.Hub
	jwtSecret       string
}

func NewDatabaseMutationProgressHandler(
	dbService *services.DatabaseService,
	progressTracker *progress.Tracker,
	hub *realtime.Hub,
) *DatabaseMutationProgressHandler {
	return &DatabaseMutationProgressHandler{
		dbService:       dbService,
		progressTracker: progressTracker,
		hub:             hub,
	}
}

func (h *DatabaseMutationProgressHandler) SetJWTSecret(secret string) {
	h.jwtSecret = secret
}

func (h *DatabaseMutationProgressHandler) CreateProjectDatabaseWithProgress(w http.ResponseWriter, r *http.Request) {

	userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		apiErr := errors.NewUnauthorizedError()
		errors.WriteError(w, apiErr)
		return
	}

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

	go func() {
		_, err := h.dbService.CreateProjectDatabaseWithProgress(r.Context(), projectID, &req, userID)
		if err != nil {

			return
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Database creation started",
		"status":  "processing",
		"note":    "Progress updates will be sent via WebSocket",
	})
}

func (h *DatabaseMutationProgressHandler) CreateTableWithProgress(w http.ResponseWriter, r *http.Request) {

	userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		apiErr := errors.NewUnauthorizedError()
		errors.WriteError(w, apiErr)
		return
	}

	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid project ID format").WithField("projectId")
		errors.WriteError(w, apiErr)
		return
	}

	var req entities.TableCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.NewValidationError("Invalid request body").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	go func() {
		_, err := h.dbService.CreateTableWithProgress(r.Context(), projectID, &req, userID)
		if err != nil {

			return
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Table creation started",
		"status":  "processing",
		"note":    "Progress updates will be sent via WebSocket",
	})
}

func (h *DatabaseMutationProgressHandler) GetOperationStatus(w http.ResponseWriter, r *http.Request) {

	userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		apiErr := errors.NewUnauthorizedError()
		errors.WriteError(w, apiErr)
		return
	}

	operationID := chi.URLParam(r, "operationId")
	operation, err := h.progressTracker.GetOperation(operationID)
	if err != nil {
		apiErr := errors.NewNotFoundError("Operation")
		errors.WriteError(w, apiErr)
		return
	}

	if operation.UserID != userID {
		apiErr := errors.NewForbiddenError()
		errors.WriteError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operation)
}

func (h *DatabaseMutationProgressHandler) GetUserOperations(w http.ResponseWriter, r *http.Request) {

	userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		apiErr := errors.NewUnauthorizedError()
		errors.WriteError(w, apiErr)
		return
	}

	operations := h.progressTracker.GetUserOperations(userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operations)
}

func (h *DatabaseMutationProgressHandler) CancelOperation(w http.ResponseWriter, r *http.Request) {

	userID, ok := authMiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		apiErr := errors.NewUnauthorizedError()
		errors.WriteError(w, apiErr)
		return
	}

	operationID := chi.URLParam(r, "operationId")
	operation, err := h.progressTracker.GetOperation(operationID)
	if err != nil {
		apiErr := errors.NewNotFoundError("Operation")
		errors.WriteError(w, apiErr)
		return
	}

	if operation.UserID != userID {
		apiErr := errors.NewForbiddenError()
		errors.WriteError(w, apiErr)
		return
	}

	if err := h.progressTracker.CancelOperation(operationID); err != nil {
		apiErr := errors.NewInternalError(err)
		errors.WriteError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Operation cancelled",
		"status":  "cancelled",
	})
}

func (h *DatabaseMutationProgressHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	realtime.HandleWebSocket(h.hub, h.jwtSecret)(w, r)
}
