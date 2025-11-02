package handlers

import (
	"net/http"
	"strconv"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/errors"
)

type TableDataHandler struct {
	tableDataService services.TableDataServiceInterface
}

func NewTableDataHandler(tableDataService services.TableDataServiceInterface) *TableDataHandler {
	return &TableDataHandler{
		tableDataService: tableDataService,
	}
}

func (h *TableDataHandler) GetTableData(w http.ResponseWriter, r *http.Request) {
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

	page := 1
	limit := 50

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	response, err := h.tableDataService.GetTableData(r.Context(), projectID, tableID, page, limit)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *TableDataHandler) GetRowByID(w http.ResponseWriter, r *http.Request) {
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

	rowID := r.URL.Query().Get("rowId")
	if rowID == "" {
		errors.WriteError(w, errors.NewValidationError("Row ID is required"))
		return
	}

	row, err := h.tableDataService.GetRowByID(r.Context(), projectID, tableID, rowID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, row)
}

func (h *TableDataHandler) InsertRow(w http.ResponseWriter, r *http.Request) {
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

	var data map[string]interface{}
	if err := decodeJSON(r, &data); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	row, err := h.tableDataService.InsertRow(r.Context(), projectID, tableID, data)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, row)
}

func (h *TableDataHandler) UpdateRow(w http.ResponseWriter, r *http.Request) {
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

	rowID, err := parseUUIDParam(r, "rowId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid row ID"))
		return
	}

	var data map[string]interface{}
	if err := decodeJSON(r, &data); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if err := h.tableDataService.UpdateRow(r.Context(), projectID, tableID, rowID, data); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *TableDataHandler) DeleteRow(w http.ResponseWriter, r *http.Request) {
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

	rowID, err := parseUUIDParam(r, "rowId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid row ID"))
		return
	}

	if err := h.tableDataService.DeleteRow(r.Context(), projectID, tableID, rowID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
