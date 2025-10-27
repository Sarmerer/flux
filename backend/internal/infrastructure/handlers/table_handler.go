package handlers

import (
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	apierrors "github.com/flow/internal/errors"
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
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	var req entities.TableCreateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid request body"))
		return
	}

	table, err := h.tableService.CreateTable(r.Context(), &req, projectID)
	if err != nil {
		if isValidationError(err) {
			writeError(w, r, apierrors.NewValidationError(err.Error()))
		} else {
			writeError(w, r, apierrors.NewInternalError(err))
		}
		return
	}

	writeJSON(w, http.StatusCreated, table)
}

func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	tableID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid table ID"))
		return
	}

	table, err := h.tableService.GetTableByID(r.Context(), tableID, projectID)
	if err != nil {
		writeError(w, r, apierrors.NewNotFoundError("Table"))
		return
	}

	writeJSON(w, http.StatusOK, table)
}

func (h *TableHandler) GetTables(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	tables, err := h.tableService.GetTablesByProjectID(r.Context(), projectID)
	if err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, tables)
}

func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	tableID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid table ID"))
		return
	}

	var req entities.TableUpdateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid request body"))
		return
	}

	table, err := h.tableService.UpdateTable(r.Context(), tableID, &req, projectID)
	if err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, table)
}

func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid project ID"))
		return
	}

	tableID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, r, apierrors.NewValidationError("Invalid table ID"))
		return
	}

	if err := h.tableService.DeleteTable(r.Context(), tableID, projectID); err != nil {
		writeError(w, r, apierrors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
