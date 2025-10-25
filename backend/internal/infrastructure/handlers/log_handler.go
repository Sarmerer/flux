package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/logging"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LogHandler struct {
	storage  logging.LogStorage
	streamer logging.LogStreamer
}

func NewLogHandler(storage logging.LogStorage, streamer logging.LogStreamer) *LogHandler {
	return &LogHandler{
		storage:  storage,
		streamer: streamer,
	}
}

func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	filter := logging.DefaultFilter()

	if projectIDStr := r.URL.Query().Get("project_id"); projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			errors.WriteError(w, errors.NewValidationError("invalid project_id").WithField("project_id"))
			return
		}
		filter.ProjectID = &projectID
	}

	if databaseIDStr := r.URL.Query().Get("database_id"); databaseIDStr != "" {
		databaseID, err := uuid.Parse(databaseIDStr)
		if err != nil {
			errors.WriteError(w, errors.NewValidationError("invalid database_id").WithField("database_id"))
			return
		}
		filter.DatabaseID = &databaseID
	}

	if tableIDStr := r.URL.Query().Get("table_id"); tableIDStr != "" {
		tableID, err := uuid.Parse(tableIDStr)
		if err != nil {
			errors.WriteError(w, errors.NewValidationError("invalid table_id").WithField("table_id"))
			return
		}
		filter.TableID = &tableID
	}

	if workflowIDStr := r.URL.Query().Get("workflow_id"); workflowIDStr != "" {
		workflowID, err := uuid.Parse(workflowIDStr)
		if err != nil {
			errors.WriteError(w, errors.NewValidationError("invalid workflow_id").WithField("workflow_id"))
			return
		}
		filter.WorkflowID = &workflowID
	}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			errors.WriteError(w, errors.NewValidationError("invalid user_id").WithField("user_id"))
			return
		}
		filter.UserID = &userID
	}

	if levelStr := r.URL.Query().Get("level"); levelStr != "" {
		level := logging.LogLevel(levelStr)
		filter.Level = &level
	}

	if operation := r.URL.Query().Get("operation"); operation != "" {
		filter.Operation = &operation
	}

	if startTimeStr := r.URL.Query().Get("start_time"); startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			errors.WriteError(w, errors.NewValidationError("invalid start_time format").WithField("start_time"))
			return
		}
		filter.StartTime = &startTime
	}

	if endTimeStr := r.URL.Query().Get("end_time"); endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			errors.WriteError(w, errors.NewValidationError("invalid end_time format").WithField("end_time"))
			return
		}
		filter.EndTime = &endTime
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			errors.WriteError(w, errors.NewValidationError("invalid limit (must be 1-1000)").WithField("limit"))
			return
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			errors.WriteError(w, errors.NewValidationError("invalid offset").WithField("offset"))
			return
		}
		filter.Offset = offset
	}

	if orderBy := r.URL.Query().Get("order_by"); orderBy != "" {
		if orderBy != "timestamp_asc" && orderBy != "timestamp_desc" {
			errors.WriteError(w, errors.NewValidationError("invalid order_by").WithField("order_by"))
			return
		}
		filter.OrderBy = orderBy
	}

	logs, err := h.storage.Query(r.Context(), filter)
	if err != nil {
		errors.WriteError(w, errors.NewDatabaseError(err).WithDetails("failed to query logs"))
		return
	}

	count, err := h.storage.Count(r.Context(), filter)
	if err != nil {
		errors.WriteError(w, errors.NewDatabaseError(err).WithDetails("failed to count logs"))
		return
	}

	response := map[string]interface{}{
		"logs":   logs,
		"total":  count,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *LogHandler) GetProjectLogs(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		errors.WriteError(w, errors.NewValidationError("invalid project ID").WithField("projectId"))
		return
	}

	filter := logging.DefaultFilter()
	filter.ProjectID = &projectID

	h.applyCommonFilters(r, filter)

	logs, err := h.storage.Query(r.Context(), filter)
	if err != nil {
		errors.WriteError(w, errors.NewDatabaseError(err).WithDetails("failed to query project logs"))
		return
	}

	count, err := h.storage.Count(r.Context(), filter)
	if err != nil {
		errors.WriteError(w, errors.NewDatabaseError(err).WithDetails("failed to count project logs"))
		return
	}

	response := map[string]interface{}{
		"logs":       logs,
		"total":      count,
		"limit":      filter.Limit,
		"offset":     filter.Offset,
		"project_id": projectID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *LogHandler) DeleteOldLogs(w http.ResponseWriter, r *http.Request) {
	daysStr := r.URL.Query().Get("days")
	if daysStr == "" {
		errors.WriteError(w, errors.NewValidationError("days parameter is required").WithField("days"))
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		errors.WriteError(w, errors.NewValidationError("invalid days parameter").WithField("days"))
		return
	}

	duration := time.Duration(days) * 24 * time.Hour
	deleted, err := h.storage.DeleteOlderThan(r.Context(), duration)
	if err != nil {
		errors.WriteError(w, errors.NewDatabaseError(err).WithDetails("failed to delete old logs"))
		return
	}

	response := map[string]interface{}{
		"deleted": deleted,
		"message": "Old logs deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *LogHandler) applyCommonFilters(r *http.Request, filter *logging.LogFilter) {
	if levelStr := r.URL.Query().Get("level"); levelStr != "" {
		level := logging.LogLevel(levelStr)
		filter.Level = &level
	}

	if operation := r.URL.Query().Get("operation"); operation != "" {
		filter.Operation = &operation
	}

	if startTimeStr := r.URL.Query().Get("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filter.StartTime = &startTime
		}
	}

	if endTimeStr := r.URL.Query().Get("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filter.EndTime = &endTime
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 1000 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if orderBy := r.URL.Query().Get("order_by"); orderBy == "timestamp_asc" || orderBy == "timestamp_desc" {
		filter.OrderBy = orderBy
	}
}
