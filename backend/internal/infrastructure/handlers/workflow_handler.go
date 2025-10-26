package handlers

import (
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	apierrors "github.com/flow/internal/errors"
)

type WorkflowHandler struct {
	workflowService services.WorkflowServiceInterface
}

func NewWorkflowHandler(workflowService services.WorkflowServiceInterface) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
	}
}

func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	var req entities.WorkflowCreateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid request body"))
		return
	}

	workflow, err := h.workflowService.CreateWorkflow(r.Context(), &req, projectID)
	if err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusCreated, workflow)
}

func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	workflowID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid workflow ID"))
		return
	}

	workflow, err := h.workflowService.GetWorkflowByID(r.Context(), workflowID, projectID)
	if err != nil {
		writeError(w, r, apierrors.NewNotFoundError("Workflow"))
		return
	}

	writeJSON(w, http.StatusOK, workflow)
}

func (h *WorkflowHandler) GetWorkflows(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	workflows, err := h.workflowService.GetWorkflowsByProjectID(r.Context(), projectID)
	if err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, workflows)
}

func (h *WorkflowHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	workflowID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid workflow ID"))
		return
	}

	var req entities.WorkflowUpdateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid request body"))
		return
	}

	workflow, err := h.workflowService.UpdateWorkflow(r.Context(), workflowID, projectID, &req)
	if err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, workflow)
}

func (h *WorkflowHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	workflowID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid workflow ID"))
		return
	}

	if err := h.workflowService.DeleteWorkflow(r.Context(), workflowID, projectID); err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkflowHandler) ToggleWorkflowActive(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	workflowID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid workflow ID"))
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid request body"))
		return
	}

	workflow, err := h.workflowService.ToggleWorkflowActive(r.Context(), workflowID, projectID, req.IsActive)
	if err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, workflow)
}

func (h *WorkflowHandler) ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	workflowID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid workflow ID"))
		return
	}

	if err := h.workflowService.ExecuteWorkflow(r.Context(), workflowID, projectID); err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Workflow executed successfully",
	})
}
