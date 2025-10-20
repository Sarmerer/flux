package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// WorkflowHandler handles workflow-related HTTP requests
type WorkflowHandler struct {
	workflowService services.WorkflowServiceInterface
}

// NewWorkflowHandler creates a new WorkflowHandler
func NewWorkflowHandler(workflowService services.WorkflowServiceInterface) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
	}
}

// CreateWorkflow handles workflow creation
func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req entities.WorkflowCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workflow, err := h.workflowService.CreateWorkflow(r.Context(), &req, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workflow)
}

// GetWorkflow handles getting a workflow by ID
func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowIDStr := chi.URLParam(r, "id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	workflow, err := h.workflowService.GetWorkflowByID(r.Context(), workflowID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// GetWorkflows handles getting all workflows for a project
func (h *WorkflowHandler) GetWorkflows(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	workflows, err := h.workflowService.GetWorkflowsByProjectID(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

// UpdateWorkflow handles updating a workflow
func (h *WorkflowHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowIDStr := chi.URLParam(r, "id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	var req entities.WorkflowUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workflow, err := h.workflowService.UpdateWorkflow(r.Context(), workflowID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// DeleteWorkflow handles deleting a workflow
func (h *WorkflowHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowIDStr := chi.URLParam(r, "id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	if err := h.workflowService.DeleteWorkflow(r.Context(), workflowID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleWorkflowActive handles toggling the active state of a workflow
func (h *WorkflowHandler) ToggleWorkflowActive(w http.ResponseWriter, r *http.Request) {
	workflowIDStr := chi.URLParam(r, "id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workflow, err := h.workflowService.ToggleWorkflowActive(r.Context(), workflowID, req.IsActive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// ExecuteWorkflow handles manually executing a workflow
func (h *WorkflowHandler) ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowIDStr := chi.URLParam(r, "id")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	if err := h.workflowService.ExecuteWorkflow(r.Context(), workflowID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Workflow executed successfully",
	})
}
