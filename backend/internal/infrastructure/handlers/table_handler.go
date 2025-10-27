package handlers

import (
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
)

type TableHandler struct {
	tableService services.TableServiceInterface
}

func NewTableHandler(tableService services.TableServiceInterface) *TableHandler {
	return &TableHandler{
		tableService: tableService,
	}
}

func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	var req entities.TableCreateRequest
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	table, err := h.tableService.CreateTable(r.Context(), &req, projectID)
	if err != nil {
		if isValidationError(err) {
			errors.WriteError(w, errors.NewValidationError(err.Error()))
		} else {
			errors.WriteError(w, errors.NewInternalError(err))
		}
		return
	}

	writeJSON(w, http.StatusCreated, table)
}

func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableID, err := parseUUIDParam(r, "id")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid table ID"))
		return
	}

	table, err := h.tableService.GetTableByID(r.Context(), tableID, projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Table not found"))
		return
	}

	writeJSON(w, http.StatusOK, table)
}

func (h *TableHandler) GetTables(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tables, err := h.tableService.GetTablesByProjectID(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, tables)
}

func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableID, err := parseUUIDParam(r, "id")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid table ID"))
		return
	}

	var req entities.TableUpdateRequest
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	table, err := h.tableService.UpdateTable(r.Context(), tableID, &req, projectID)
	if err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, table)
}

func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableID, err := parseUUIDParam(r, "id")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid table ID"))
		return
	}

	if err := h.tableService.DeleteTable(r.Context(), tableID, projectID); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
