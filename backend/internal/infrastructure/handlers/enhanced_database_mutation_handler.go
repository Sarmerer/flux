package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	authMiddleware "github.com/flow/internal/infrastructure/middleware"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/progress"
	"github.com/flow/internal/infrastructure/realtime"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// DatabaseMutationProgressHandler handles database mutation requests with progress tracking and WebSocket support
type DatabaseMutationProgressHandler struct {
	dbService       *services.DatabaseService
	progressTracker *progress.Tracker
	hub             *realtime.Hub
	jwtSecret       string
}

// NewDatabaseMutationProgressHandler creates a new database mutation progress handler
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

// SetJWTSecret sets the JWT secret for WebSocket authentication
func (h *DatabaseMutationProgressHandler) SetJWTSecret(secret string) {
	h.jwtSecret = secret
}

// CreateProjectDatabaseWithProgress handles project database creation with progress tracking
func (h *DatabaseMutationProgressHandler) CreateProjectDatabaseWithProgress(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
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

	// Start the operation asynchronously
	go func() {
		_, err := h.dbService.CreateProjectDatabaseWithProgress(r.Context(), projectID, &req, userID)
		if err != nil {
			// Error handling is done in the service via progress tracking
			return
		}
	}()

	// Return immediately with operation started response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Database creation started",
		"status":  "processing",
		"note":    "Progress updates will be sent via WebSocket",
	})
}

// CreateTableWithProgress handles table creation with progress tracking
func (h *DatabaseMutationProgressHandler) CreateTableWithProgress(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
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

	// Start the operation asynchronously
	go func() {
		_, err := h.dbService.CreateTableWithProgress(r.Context(), projectID, &req, userID)
		if err != nil {
			// Error handling is done in the service via progress tracking
			return
		}
	}()

	// Return immediately with operation started response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Table creation started",
		"status":  "processing",
		"note":    "Progress updates will be sent via WebSocket",
	})
}

// GetOperationStatus returns the status of a specific operation
func (h *DatabaseMutationProgressHandler) GetOperationStatus(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
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

	// Check if user owns this operation
	if operation.UserID != userID {
		apiErr := errors.NewForbiddenError()
		errors.WriteError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operation)
}

// GetUserOperations returns all operations for the current user
func (h *DatabaseMutationProgressHandler) GetUserOperations(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
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

// CancelOperation cancels a running operation
func (h *DatabaseMutationProgressHandler) CancelOperation(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
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

	// Check if user owns this operation
	if operation.UserID != userID {
		apiErr := errors.NewForbiddenError()
		errors.WriteError(w, apiErr)
		return
	}

	// Cancel the operation
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

// WebSocketHandler handles WebSocket connections
func (h *DatabaseMutationProgressHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	realtime.HandleWebSocket(h.hub, h.jwtSecret)(w, r)
}
