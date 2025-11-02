package handlers

import (
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/validation"
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
		writeError(w, err)
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
		writeError(w, err)
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
		writeError(w, err)
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
		writeError(w, err)
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
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TableHandler) AddColumn(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Column database.ColumnDefinition `json:"column"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if err := h.tableService.AddColumn(r.Context(), tableID, projectID, req.Column); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "column added"})
}

func (h *TableHandler) RemoveColumn(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		ColumnName string `json:"column_name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if err := h.tableService.RemoveColumn(r.Context(), tableID, projectID, req.ColumnName); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TableHandler) ModifyColumn(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		ColumnName string                   `json:"column_name"`
		Column     database.ColumnDefinition `json:"column"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if err := h.tableService.ModifyColumn(r.Context(), tableID, projectID, req.ColumnName, req.Column); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "column modified"})
}

func (h *TableHandler) AddForeignKey(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		ForeignKey database.ForeignKeyDefinition `json:"foreign_key"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if err := h.tableService.AddForeignKey(r.Context(), tableID, projectID, req.ForeignKey); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "foreign key added"})
}

func (h *TableHandler) RemoveForeignKey(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		ForeignKeyName string `json:"foreign_key_name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if err := h.tableService.RemoveForeignKey(r.Context(), tableID, projectID, req.ForeignKeyName); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TableHandler) UpdateTableSchema(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Schema database.TableSchema `json:"schema"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if err := validation.ValidateTableSchema(&req.Schema); err != nil {
		errors.WriteError(w, errors.NewValidationError("Schema validation failed").WithDetails(err.Error()))
		return
	}

	if err := h.tableService.UpdateTableSchema(r.Context(), tableID, projectID, req.Schema); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "schema updated"})
}
