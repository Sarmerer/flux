package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/pkg/errors"
	"github.com/flow/internal/pkg/progress"
	"github.com/flow/internal/pkg/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// EnhancedDatabaseMutationHandler handles database mutation requests with progress tracking and WebSocket support
type EnhancedDatabaseMutationHandler struct {
	enhancedService *services.EnhancedDatabaseMutationService
	progressTracker *progress.Tracker
	hub             *websocket.Hub
}

// NewEnhancedDatabaseMutationHandler creates a new enhanced database mutation handler
func NewEnhancedDatabaseMutationHandler(
	enhancedService *services.EnhancedDatabaseMutationService,
	progressTracker *progress.Tracker,
	hub *websocket.Hub,
) *EnhancedDatabaseMutationHandler {
	return &EnhancedDatabaseMutationHandler{
		enhancedService: enhancedService,
		progressTracker: progressTracker,
		hub:             hub,
	}
}

// CreateProjectDatabaseWithProgress handles project database creation with progress tracking
func (h *EnhancedDatabaseMutationHandler) CreateProjectDatabaseWithProgress(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
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
		_, err := h.enhancedService.CreateProjectDatabaseWithProgress(r.Context(), projectID, &req, userID)
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
func (h *EnhancedDatabaseMutationHandler) CreateTableWithProgress(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
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
		_, err := h.enhancedService.CreateTableWithProgress(r.Context(), projectID, &req, userID)
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
func (h *EnhancedDatabaseMutationHandler) GetOperationStatus(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
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
func (h *EnhancedDatabaseMutationHandler) GetUserOperations(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
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
func (h *EnhancedDatabaseMutationHandler) CancelOperation(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
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
func (h *EnhancedDatabaseMutationHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	websocket.HandleWebSocket(h.hub)(w, r)
}
