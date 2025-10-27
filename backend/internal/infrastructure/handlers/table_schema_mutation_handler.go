package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/validation"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TableSchemaMutationHandler struct {
	dbRepo      repositories.DatabaseRepository
	projectRepo repositories.ProjectRepository
	schemaSvc   *database.SchemaManagementService
}

type ColumnDefinition = database.ColumnDefinition
type TableSchema = database.TableSchema
type ForeignKeyDefinition = database.ForeignKeyDefinition

func NewTableSchemaMutationHandler(
	dbRepo repositories.DatabaseRepository,
	projectRepo repositories.ProjectRepository,
	schemaSvc *database.SchemaManagementService,
) *TableSchemaMutationHandler {
	return &TableSchemaMutationHandler{
		dbRepo:      dbRepo,
		projectRepo: projectRepo,
		schemaSvc:   schemaSvc,
	}
}

func (h *TableSchemaMutationHandler) getProjectDatabase(ctx context.Context, projectID uuid.UUID) (*entities.Database, error) {
	databases, err := h.dbRepo.GetByProjectID(ctx, projectID)
	if err != nil || len(databases) == 0 {
		return nil, fmt.Errorf("no database found for project")
	}
	return databases[0], nil
}

func (h *TableSchemaMutationHandler) CreateTableInDatabase(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	var req struct {
		TableName string      `json:"table_name" validate:"required"`
		Schema    TableSchema `json:"schema" validate:"required"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	if req.TableName == "" {
		errors.WriteError(w, errors.NewValidationError("Table name is required"))
		return
	}

	if !validation.IsValidTableName(req.TableName) {
		errors.WriteError(w, errors.NewValidationError("Invalid table name"))
		return
	}

	if err := validation.ValidateTableSchema(&req.Schema); err != nil {
		errors.WriteError(w, errors.NewValidationError("Schema validation failed").WithDetails(err.Error()))
		return
	}

	database, err := h.getProjectDatabase(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project database not found"))
		return
	}

	if err := h.schemaSvc.CreateTable(r.Context(), database, req.TableName, req.Schema); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "table created"})
}

func (h *TableSchemaMutationHandler) DropTableFromDatabase(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		errors.WriteError(w, errors.NewValidationError("Table name is required"))
		return
	}

	database, err := h.getProjectDatabase(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project database not found"))
		return
	}

	if err := h.schemaSvc.DropTable(r.Context(), database, tableName); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TableSchemaMutationHandler) AddColumnToTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		errors.WriteError(w, errors.NewValidationError("Table name is required"))
		return
	}

	var req struct {
		Column ColumnDefinition `json:"column" validate:"required"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	database, err := h.getProjectDatabase(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project database not found"))
		return
	}

	if err := h.schemaSvc.AddColumn(r.Context(), database, tableName, req.Column); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "column added"})
}

func (h *TableSchemaMutationHandler) RemoveColumnFromTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		errors.WriteError(w, errors.NewValidationError("Table name is required"))
		return
	}

	var req struct {
		ColumnName string `json:"column_name" validate:"required"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	database, err := h.getProjectDatabase(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project database not found"))
		return
	}

	if err := h.schemaSvc.RemoveColumn(r.Context(), database, tableName, req.ColumnName); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TableSchemaMutationHandler) ModifyColumnInTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		errors.WriteError(w, errors.NewValidationError("Table name is required"))
		return
	}

	var req struct {
		ColumnName string           `json:"column_name" validate:"required"`
		Column     ColumnDefinition `json:"column" validate:"required"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	database, err := h.getProjectDatabase(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project database not found"))
		return
	}

	if err := h.schemaSvc.ModifyColumn(r.Context(), database, tableName, req.ColumnName, req.Column); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "column modified"})
}

func (h *TableSchemaMutationHandler) AddForeignKeyToTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		errors.WriteError(w, errors.NewValidationError("Table name is required"))
		return
	}

	var req struct {
		ForeignKey ForeignKeyDefinition `json:"foreign_key" validate:"required"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	database, err := h.getProjectDatabase(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project database not found"))
		return
	}

	if err := h.schemaSvc.AddForeignKey(r.Context(), database, tableName, req.ForeignKey); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "foreign key added"})
}

func (h *TableSchemaMutationHandler) RemoveForeignKeyFromTable(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectId")
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid project ID"))
		return
	}

	tableName := chi.URLParam(r, "tableName")
	if tableName == "" {
		errors.WriteError(w, errors.NewValidationError("Table name is required"))
		return
	}

	var req struct {
		ForeignKeyName string `json:"foreign_key_name" validate:"required"`
	}
	if err := decodeJSON(r, &req); err != nil {
		errors.WriteError(w, errors.NewValidationError("Invalid request body"))
		return
	}

	database, err := h.getProjectDatabase(r.Context(), projectID)
	if err != nil {
		errors.WriteError(w, errors.NewNotFoundError("Project database not found"))
		return
	}

	if err := h.schemaSvc.RemoveForeignKey(r.Context(), database, tableName, req.ForeignKeyName); err != nil {
		errors.WriteError(w, errors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
