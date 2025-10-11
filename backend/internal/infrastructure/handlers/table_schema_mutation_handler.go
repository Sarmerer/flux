package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/app/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TableSchemaMutationHandler handles table schema mutation HTTP requests
type TableSchemaMutationHandler struct {
	schemaMutationService *services.TableSchemaMutationService
}

// NewTableSchemaMutationHandler creates a new TableSchemaMutationHandler
func NewTableSchemaMutationHandler(schemaMutationService *services.TableSchemaMutationService) *TableSchemaMutationHandler {
	return &TableSchemaMutationHandler{
		schemaMutationService: schemaMutationService,
	}
}

// CreateTableInDatabase handles creating a table in the project database
func (h *TableSchemaMutationHandler) CreateTableInDatabase(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req struct {
		TableName string               `json:"table_name" validate:"required"`
		Schema    services.TableSchema `json:"schema" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.schemaMutationService.CreateTableInDatabase(r.Context(), projectID, req.TableName, req.Schema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "table created"})
}

// DropTableFromDatabase handles dropping a table from the project database
func (h *TableSchemaMutationHandler) DropTableFromDatabase(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	if err := h.schemaMutationService.DropTableFromDatabase(r.Context(), projectID, tableName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddColumnToTable handles adding a column to a table
func (h *TableSchemaMutationHandler) AddColumnToTable(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	var req services.AddColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.schemaMutationService.AddColumnToTable(r.Context(), projectID, tableName, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "column added"})
}

// RemoveColumnFromTable handles removing a column from a table
func (h *TableSchemaMutationHandler) RemoveColumnFromTable(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	var req services.RemoveColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.schemaMutationService.RemoveColumnFromTable(r.Context(), projectID, tableName, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ModifyColumnInTable handles modifying a column in a table
func (h *TableSchemaMutationHandler) ModifyColumnInTable(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	var req services.ModifyColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.schemaMutationService.ModifyColumnInTable(r.Context(), projectID, tableName, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "column modified"})
}

// AddForeignKeyToTable handles adding a foreign key to a table
func (h *TableSchemaMutationHandler) AddForeignKeyToTable(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	var req services.AddForeignKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.schemaMutationService.AddForeignKeyToTable(r.Context(), projectID, tableName, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "foreign key added"})
}

// RemoveForeignKeyFromTable handles removing a foreign key from a table
func (h *TableSchemaMutationHandler) RemoveForeignKeyFromTable(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	var req services.RemoveForeignKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.schemaMutationService.RemoveForeignKeyFromTable(r.Context(), projectID, tableName, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
