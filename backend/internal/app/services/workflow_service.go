package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/flow/internal/errors"

	"github.com/google/uuid"
)

// WorkflowService handles workflow-related business logic
type WorkflowService struct {
	workflowRepo repositories.WorkflowRepository
	projectRepo  repositories.ProjectRepository
	dbRepo       repositories.DatabaseRepository
	pgService    *database.PostgreSQLManagementService
	httpClient   *http.Client
	logger       *logging.Logger
}

// NewWorkflowService creates a new WorkflowService
func NewWorkflowService(
	workflowRepo repositories.WorkflowRepository,
	projectRepo repositories.ProjectRepository,
	dbRepo repositories.DatabaseRepository,
	pgService *database.PostgreSQLManagementService,
	logger *logging.Logger,
) *WorkflowService {
	return &WorkflowService{
		workflowRepo: workflowRepo,
		projectRepo:  projectRepo,
		dbRepo:       dbRepo,
		pgService:    pgService,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CreateWorkflow creates a new workflow
func (s *WorkflowService) CreateWorkflow(ctx context.Context, req *entities.WorkflowCreateRequest, projectID uuid.UUID) (*entities.WorkflowResponse, error) {
	// Verify project exists
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewNotFoundError("Project")
	}

	// Validate trigger
	if err := s.validateTrigger(&req.Trigger); err != nil {
		return nil, errors.NewValidationError("invalid trigger").WithDetails(err.Error())
	}

	// Validate actions
	if err := s.validateActions(req.Actions); err != nil {
		return nil, errors.NewValidationError("invalid actions").WithDetails(err.Error())
	}

	// Create workflow
	workflow := &entities.Workflow{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Trigger:     req.Trigger,
		Actions:     req.Actions,
		IsActive:    false, // Workflows start inactive by default
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.workflowRepo.Create(ctx, workflow); err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to create workflow")
	}

	response := workflow.ToResponse()
	return &response, nil
}

// GetWorkflowByID retrieves a workflow by ID
func (s *WorkflowService) GetWorkflowByID(ctx context.Context, id uuid.UUID) (*entities.WorkflowResponse, error) {
	workflow, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Workflow")
	}

	response := workflow.ToResponse()
	return &response, nil
}

// GetWorkflowsByProjectID retrieves all workflows for a project
func (s *WorkflowService) GetWorkflowsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.WorkflowResponse, error) {
	workflows, err := s.workflowRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to get workflows")
	}

	var responses []*entities.WorkflowResponse
	for _, workflow := range workflows {
		response := workflow.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

// UpdateWorkflow updates a workflow
func (s *WorkflowService) UpdateWorkflow(ctx context.Context, id uuid.UUID, req *entities.WorkflowUpdateRequest) (*entities.WorkflowResponse, error) {
	workflow, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Workflow")
	}

	// Update fields if provided
	if req.Name != "" {
		workflow.Name = req.Name
	}

	if req.Description != "" {
		workflow.Description = req.Description
	}

	if req.Trigger != nil {
		if err := s.validateTrigger(req.Trigger); err != nil {
			return nil, errors.NewValidationError("invalid trigger").WithDetails(err.Error())
		}
		workflow.Trigger = *req.Trigger
	}

	if req.Actions != nil && len(req.Actions) > 0 {
		if err := s.validateActions(req.Actions); err != nil {
			return nil, errors.NewValidationError("invalid actions").WithDetails(err.Error())
		}
		workflow.Actions = req.Actions
	}

	if req.IsActive != nil {
		workflow.IsActive = *req.IsActive
	}

	workflow.UpdatedAt = time.Now()

	if err := s.workflowRepo.Update(ctx, workflow); err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to update workflow")
	}

	response := workflow.ToResponse()
	return &response, nil
}

// DeleteWorkflow deletes a workflow
func (s *WorkflowService) DeleteWorkflow(ctx context.Context, id uuid.UUID) error {
	// Check if workflow exists
	_, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Workflow")
	}

	if err := s.workflowRepo.Delete(ctx, id); err != nil {
		return errors.NewDatabaseError(err).WithDetails("failed to delete workflow")
	}

	return nil
}

// ToggleWorkflowActive toggles the active state of a workflow
func (s *WorkflowService) ToggleWorkflowActive(ctx context.Context, id uuid.UUID, isActive bool) (*entities.WorkflowResponse, error) {
	workflow, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Workflow")
	}

	if err := s.workflowRepo.ToggleActive(ctx, id, isActive); err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to toggle workflow active state")
	}

	// Fetch updated workflow
	workflow, err = s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError(err).WithDetails("failed to fetch updated workflow")
	}

	response := workflow.ToResponse()
	return &response, nil
}

// ExecuteWorkflow manually executes a workflow for testing
func (s *WorkflowService) ExecuteWorkflow(ctx context.Context, id uuid.UUID) error {
	workflow, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Workflow")
	}

	// Create context logger for workflow execution
	logger := s.logger.WithContext(logging.LogContext{
		WorkflowID: &workflow.ID,
		ProjectID:  &workflow.ProjectID,
		Operation:  "execute_workflow",
	})

	logger.Info("Executing workflow", map[string]interface{}{
		"workflow_name": workflow.Name,
		"action_count":  len(workflow.Actions),
	})

	startTime := time.Now()

	// Execute each action in order
	for i, action := range workflow.Actions {
		actionLogger := s.logger.WithContext(logging.LogContext{
			WorkflowID: &workflow.ID,
			ProjectID:  &workflow.ProjectID,
			Operation:  fmt.Sprintf("execute_action_%s", action.Type),
		})

		actionLogger.Info("Executing workflow action", map[string]interface{}{
			"action_index":  i + 1,
			"total_actions": len(workflow.Actions),
			"action_id":     action.ID,
			"action_type":   action.Type,
		})

		actionStart := time.Now()

		if err := s.executeAction(ctx, workflow, &action); err != nil {
			actionLogger.Error("Failed to execute workflow action", err, map[string]interface{}{
				"action_id":   action.ID,
				"action_type": action.Type,
			})
			return errors.NewAPIError(errors.ErrCodeOperationFailed, fmt.Sprintf("failed to execute action %s", action.ID)).WithDetails(err.Error())
		}

		actionLogger.Info("Workflow action completed", map[string]interface{}{
			"action_id":   action.ID,
			"action_type": action.Type,
			"duration_ms": time.Since(actionStart).Milliseconds(),
		})
	}

	logger.Info("Workflow execution completed", map[string]interface{}{
		"workflow_name": workflow.Name,
		"action_count":  len(workflow.Actions),
		"duration_ms":   time.Since(startTime).Milliseconds(),
	})

	return nil
}

// validateTrigger validates a workflow trigger
func (s *WorkflowService) validateTrigger(trigger *entities.WorkflowTrigger) error {
	validTriggerTypes := map[string]bool{
		"on_row_created": true,
		"on_row_updated": true,
		"on_row_deleted": true,
		"scheduled":      true,
		"webhook":        true,
	}

	if !validTriggerTypes[trigger.Type] {
		return fmt.Errorf("invalid trigger type: %s", trigger.Type)
	}

	// Validate table-based triggers
	if trigger.Type == "on_row_created" || trigger.Type == "on_row_updated" || trigger.Type == "on_row_deleted" {
		if trigger.TableName == "" {
			return fmt.Errorf("table_name is required for trigger type: %s", trigger.Type)
		}
	}

	// Validate scheduled triggers with cron expression
	if trigger.Type == "scheduled" {
		if trigger.Schedule == "" {
			return fmt.Errorf("schedule is required for scheduled triggers")
		}
		if err := validateCronExpression(trigger.Schedule); err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}
	}

	return nil
}

// validateCronExpression validates a cron expression format
// Supports standard 5-field cron format: minute hour day month weekday
func validateCronExpression(expr string) error {
	// Basic cron expression validation using regex
	// Format: "minute hour day month weekday"
	// Each field can be: *, number, range (1-5), list (1,2,3), or step (*/5)
	cronRegex := regexp.MustCompile(`^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|([0-9]|1[0-9]|2[0-3])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|([1-9]|1[0-2])|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+(\*|[0-6]|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)$`)

	if !cronRegex.MatchString(expr) {
		return fmt.Errorf("invalid cron expression format, expected: 'minute hour day month weekday'")
	}

	return nil
}

// validateActions validates workflow actions
func (s *WorkflowService) validateActions(actions []entities.WorkflowAction) error {
	if len(actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}

	validActionTypes := map[string]bool{
		"send_webhook": true,
		"send_email":   true,
		"update_row":   true,
		"create_row":   true,
		"delete_row":   true,
	}

	for i, action := range actions {
		if action.ID == "" {
			return fmt.Errorf("action %d: id is required", i)
		}

		if !validActionTypes[action.Type] {
			return fmt.Errorf("action %d: invalid action type: %s", i, action.Type)
		}

		if action.Config == nil || len(action.Config) == 0 {
			return fmt.Errorf("action %d: config is required", i)
		}

		// Validate action-specific requirements
		if err := s.validateActionConfig(i, &action); err != nil {
			return err
		}
	}

	return nil
}

// validateActionConfig validates action-specific configuration
func (s *WorkflowService) validateActionConfig(index int, action *entities.WorkflowAction) error {
	switch action.Type {
	case "send_webhook":
		if _, ok := action.Config["url"].(string); !ok {
			return fmt.Errorf("action %d: webhook url is required", index)
		}
	case "send_email":
		if _, ok := action.Config["to"].(string); !ok {
			return fmt.Errorf("action %d: email recipient (to) is required", index)
		}
	case "update_row", "create_row", "delete_row":
		if _, ok := action.Config["table"].(string); !ok {
			return fmt.Errorf("action %d: table name is required", index)
		}
	}
	return nil
}

// executeAction executes a single workflow action
func (s *WorkflowService) executeAction(ctx context.Context, workflow *entities.Workflow, action *entities.WorkflowAction) error {
	switch action.Type {
	case "send_webhook":
		return s.executeSendWebhook(ctx, action)
	case "send_email":
		return s.executeSendEmail(ctx, action)
	case "update_row":
		return s.executeUpdateRow(ctx, workflow, action)
	case "create_row":
		return s.executeCreateRow(ctx, workflow, action)
	case "delete_row":
		return s.executeDeleteRow(ctx, workflow, action)
	default:
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// executeSendWebhook executes a send_webhook action
func (s *WorkflowService) executeSendWebhook(ctx context.Context, action *entities.WorkflowAction) error {
	url, ok := action.Config["url"].(string)
	if !ok {
		return fmt.Errorf("webhook url not found in config")
	}

	// Prepare webhook payload
	payload := map[string]interface{}{
		"action_id": action.ID,
		"type":      action.Type,
		"timestamp": time.Now().Unix(),
		"data":      action.Config,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Flow-Workflow-Engine/1.0")

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d", resp.StatusCode)
	}

	s.logger.Info("Successfully sent webhook", map[string]interface{}{
		"url":    url,
		"status": resp.StatusCode,
	})
	return nil
}

// executeSendEmail executes a send_email action
func (s *WorkflowService) executeSendEmail(ctx context.Context, action *entities.WorkflowAction) error {
	to, ok := action.Config["to"].(string)
	if !ok {
		return fmt.Errorf("email recipient not found in config")
	}

	subject, _ := action.Config["subject"].(string)
	body, _ := action.Config["body"].(string)

	// For now, just log the email that would be sent
	// In production, integrate with an email service (SendGrid, SES, SMTP, etc.)
	s.logger.Info("Email would be sent (email service not configured)", map[string]interface{}{
		"to":      to,
		"subject": subject,
		"body":    body,
	})

	// TODO: Integrate with actual email service when ready
	// Example integration points:
	// - AWS SES: use aws-sdk-go-v2/service/ses
	// - SendGrid: use sendgrid-go library
	// - SMTP: use net/smtp package

	return nil
}

// executeUpdateRow executes an update_row action
func (s *WorkflowService) executeUpdateRow(ctx context.Context, workflow *entities.Workflow, action *entities.WorkflowAction) error {
	tableName, ok := action.Config["table"].(string)
	if !ok {
		return fmt.Errorf("table name not found in config")
	}

	rowID, ok := action.Config["row_id"].(string)
	if !ok {
		return fmt.Errorf("row_id not found in config")
	}

	updates, ok := action.Config["updates"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("updates not found in config")
	}

	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, workflow.ProjectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("project database not found")
	}

	database := databases[0]

	// Build UPDATE query dynamically
	// Note: In production, use parameterized queries to prevent SQL injection
	setClause := ""
	args := []interface{}{}
	argPos := 1

	for key, value := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", key, argPos)
		args = append(args, value)
		argPos++
	}

	args = append(args, rowID)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", tableName, setClause, argPos)

	// Execute query
	return s.pgService.ExecuteQuery(ctx, database, query, args...)
}

// executeCreateRow executes a create_row action
func (s *WorkflowService) executeCreateRow(ctx context.Context, workflow *entities.Workflow, action *entities.WorkflowAction) error {
	tableName, ok := action.Config["table"].(string)
	if !ok {
		return fmt.Errorf("table name not found in config")
	}

	data, ok := action.Config["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("data not found in config")
	}

	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, workflow.ProjectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("project database not found")
	}

	database := databases[0]

	// Build INSERT query dynamically
	columns := ""
	placeholders := ""
	args := []interface{}{}
	argPos := 1

	for key, value := range data {
		if columns != "" {
			columns += ", "
			placeholders += ", "
		}
		columns += key
		placeholders += fmt.Sprintf("$%d", argPos)
		args = append(args, value)
		argPos++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)

	// Execute query
	return s.pgService.ExecuteQuery(ctx, database, query, args...)
}

// executeDeleteRow executes a delete_row action
func (s *WorkflowService) executeDeleteRow(ctx context.Context, workflow *entities.Workflow, action *entities.WorkflowAction) error {
	tableName, ok := action.Config["table"].(string)
	if !ok {
		return fmt.Errorf("table name not found in config")
	}

	rowID, ok := action.Config["row_id"].(string)
	if !ok {
		return fmt.Errorf("row_id not found in config")
	}

	// Get project database
	databases, err := s.dbRepo.GetByProjectID(ctx, workflow.ProjectID)
	if err != nil || len(databases) == 0 {
		return fmt.Errorf("project database not found")
	}

	database := databases[0]

	// Build DELETE query
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	// Execute query
	return s.pgService.ExecuteQuery(ctx, database, query, rowID)
}
