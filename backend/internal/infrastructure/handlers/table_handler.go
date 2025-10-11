package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TableHandler handles table-related HTTP requests
type TableHandler struct {
	tableService services.TableServiceInterface
}

// NewTableHandler creates a new TableHandler
func NewTableHandler(tableService services.TableServiceInterface) *TableHandler {
	return &TableHandler{
		tableService: tableService,
	}
}

// CreateTable handles table creation
func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req entities.TableCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	table, err := h.tableService.CreateTable(r.Context(), &req, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(table)
}

// GetTable handles getting a table by ID
func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	tableIDStr := chi.URLParam(r, "id")
	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	table, err := h.tableService.GetTableByID(r.Context(), tableID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

// GetTables handles getting all tables for a project
func (h *TableHandler) GetTables(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	tables, err := h.tableService.GetTablesByProjectID(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}

// UpdateTable handles updating a table
func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	tableIDStr := chi.URLParam(r, "id")
	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	var req entities.TableUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	table, err := h.tableService.UpdateTable(r.Context(), tableID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

// DeleteTable handles deleting a table
func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	tableIDStr := chi.URLParam(r, "id")
	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	if err := h.tableService.DeleteTable(r.Context(), tableID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
