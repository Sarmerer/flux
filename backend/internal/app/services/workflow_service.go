package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/logging"
	"github.com/flow/internal/validation"

	"github.com/google/uuid"
)

type WorkflowService struct {
	repoFactory         *ProjectRepositoryFactory
	projectRepo         repositories.ProjectRepository
	dbRepo              repositories.DatabaseRepository
	dataManipulationSvc *database.DataManipulationService
	httpClient          *http.Client
	logger              *logging.Logger
}

func NewWorkflowService(
	repoFactory *ProjectRepositoryFactory,
	projectRepo repositories.ProjectRepository,
	dbRepo repositories.DatabaseRepository,
	dataManipulationSvc *database.DataManipulationService,
	logger *logging.Logger,
) *WorkflowService {
	return &WorkflowService{
		repoFactory:         repoFactory,
		projectRepo:         projectRepo,
		dbRepo:              dbRepo,
		dataManipulationSvc: dataManipulationSvc,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (s *WorkflowService) CreateWorkflow(ctx context.Context, req *entities.WorkflowCreateRequest, projectID uuid.UUID) (*entities.WorkflowResponse, error) {
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, errors.NewNotFoundError("Project not found")
	}

	if err := validation.ValidateTrigger(&req.Trigger); err != nil {
		return nil, errors.NewValidationError("invalid trigger").WithDetails(err.Error())
	}

	if err := validation.ValidateActions(req.Actions); err != nil {
		return nil, errors.NewValidationError("invalid actions").WithDetails(err.Error())
	}

	workflowRepo, err := s.repoFactory.GetWorkflowRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get workflow repository", err)
	}

	workflow := &entities.Workflow{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Trigger:     req.Trigger,
		Actions:     req.Actions,
		IsActive:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := workflowRepo.Create(ctx, workflow); err != nil {
		return nil, errors.NewDatabaseError("Failed to create workflow", err)
	}

	response := workflow.ToResponse()
	return &response, nil
}

func (s *WorkflowService) GetWorkflowByID(ctx context.Context, id uuid.UUID, projectID uuid.UUID) (*entities.WorkflowResponse, error) {
	workflowRepo, err := s.repoFactory.GetWorkflowRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get workflow repository", err)
	}

	workflow, err := workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Workflow not found")
	}

	response := workflow.ToResponse()
	return &response, nil
}

func (s *WorkflowService) GetWorkflowsByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.WorkflowResponse, error) {
	workflowRepo, err := s.repoFactory.GetWorkflowRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get workflow repository", err)
	}

	workflows, err := workflowRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get workflows", err)
	}

	responses := make([]*entities.WorkflowResponse, 0)
	for _, workflow := range workflows {
		response := workflow.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *WorkflowService) UpdateWorkflow(ctx context.Context, id uuid.UUID, projectID uuid.UUID, req *entities.WorkflowUpdateRequest) (*entities.WorkflowResponse, error) {
	workflowRepo, err := s.repoFactory.GetWorkflowRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get workflow repository", err)
	}

	workflow, err := workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Workflow not found")
	}

	if req.Name != "" {
		workflow.Name = req.Name
	}

	if req.Description != "" {
		workflow.Description = req.Description
	}

	if req.Trigger != nil {
		if err := validation.ValidateTrigger(req.Trigger); err != nil {
			return nil, errors.NewValidationError("invalid trigger").WithDetails(err.Error())
		}
		workflow.Trigger = *req.Trigger
	}

	if req.Actions != nil && len(req.Actions) > 0 {
		if err := validation.ValidateActions(req.Actions); err != nil {
			return nil, errors.NewValidationError("invalid actions").WithDetails(err.Error())
		}
		workflow.Actions = req.Actions
	}

	if req.IsActive != nil {
		workflow.IsActive = *req.IsActive
	}

	workflow.UpdatedAt = time.Now()

	if err := workflowRepo.Update(ctx, workflow); err != nil {
		return nil, errors.NewDatabaseError("Failed to update workflow", err)
	}

	response := workflow.ToResponse()
	return &response, nil
}

func (s *WorkflowService) DeleteWorkflow(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error {
	workflowRepo, err := s.repoFactory.GetWorkflowRepository(ctx, projectID)
	if err != nil {
		return errors.NewDatabaseError("Failed to get workflow repository", err)
	}

	_, err = workflowRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Workflow not found")
	}

	if err := workflowRepo.Delete(ctx, id); err != nil {
		return errors.NewDatabaseError("Failed to delete workflow", err)
	}

	return nil
}

func (s *WorkflowService) ToggleWorkflowActive(ctx context.Context, id uuid.UUID, projectID uuid.UUID, isActive bool) (*entities.WorkflowResponse, error) {
	workflowRepo, err := s.repoFactory.GetWorkflowRepository(ctx, projectID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get workflow repository", err)
	}

	workflow, err := workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Workflow not found")
	}

	if err := workflowRepo.ToggleActive(ctx, id, isActive); err != nil {
		return nil, errors.NewDatabaseError("Failed to toggle workflow active status", err)
	}

	workflow, err = workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to retrieve updated workflow", err)
	}

	response := workflow.ToResponse()
	return &response, nil
}

func (s *WorkflowService) ExecuteWorkflow(ctx context.Context, id uuid.UUID, projectID uuid.UUID) error {
	workflowRepo, err := s.repoFactory.GetWorkflowRepository(ctx, projectID)
	if err != nil {
		return errors.NewDatabaseError("Failed to get workflow repository", err)
	}

	workflow, err := workflowRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("Workflow not found")
	}

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
			return errors.NewAPIError(errors.ErrCodeOperationFailed, fmt.Sprintf("Failed to execute action %s", action.ID)).WithDetails(err.Error())
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
		return errors.NewValidationError("Unsupported action type").
			WithField("action.type").
			WithDetails(fmt.Sprintf("Type '%s' is not supported", action.Type))
	}
}

func (s *WorkflowService) executeSendWebhook(ctx context.Context, action *entities.WorkflowAction) error {
	webhookURL, ok := action.Config["url"].(string)
	if !ok {
		return errors.NewValidationError("Webhook URL not found in config").
			WithField("action.config.url")
	}

	if err := s.validateWebhookURL(webhookURL); err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return apiErr
		}
		return errors.NewValidationError("Invalid webhook URL").WithDetails(err.Error())
	}

	payload := map[string]interface{}{
		"action_id": action.ID,
		"type":      action.Type,
		"timestamp": time.Now().Unix(),
		"data":      action.Config,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to marshal webhook payload")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errors.NewInternalError(err).WithDetails("Failed to create webhook request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Flow-Workflow-Engine/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.NewAPIError(errors.ErrCodeExternalService, "Failed to send webhook").WithDetails(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.NewAPIError(errors.ErrCodeExternalService, "Webhook request failed").
			WithDetails(fmt.Sprintf("Received status code %d", resp.StatusCode))
	}

	s.logger.Info("Successfully sent webhook", map[string]interface{}{
		"url":    webhookURL,
		"status": resp.StatusCode,
	})
	return nil
}

func (s *WorkflowService) validateWebhookURL(webhookURL string) error {
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return errors.NewValidationError("Invalid URL format").
			WithField("url").
			WithDetails(err.Error())
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.NewValidationError("Only HTTP and HTTPS schemes are allowed").
			WithField("url").
			WithDetails(fmt.Sprintf("Scheme '%s' is not supported", parsedURL.Scheme))
	}

	if parsedURL.Host == "" {
		return errors.NewValidationError("URL must have a host").WithField("url")
	}

	host := parsedURL.Hostname()

	blockedHosts := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
		"[::1]",
	}

	for _, blocked := range blockedHosts {
		if strings.EqualFold(host, blocked) {
			return errors.NewValidationError("Localhost addresses are not allowed").
				WithField("url").
				WithDetails(fmt.Sprintf("Host '%s' is blocked", host))
		}
	}

	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsLoopback() {
			return errors.NewValidationError("Loopback addresses are not allowed").WithField("url")
		}
		if ip.IsPrivate() {
			return errors.NewValidationError("Private IP addresses are not allowed").WithField("url")
		}
		if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return errors.NewValidationError("Link-local addresses are not allowed").WithField("url")
		}
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return errors.NewAPIError(errors.ErrCodeExternalService, "Failed to resolve host").
			WithField("url").
			WithDetails(err.Error())
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
			return errors.NewValidationError("Host resolves to a blocked IP address").
				WithField("url").
				WithDetails(fmt.Sprintf("Host '%s' resolves to blocked IP", host))
		}
	}

	return nil
}

func (s *WorkflowService) executeSendEmail(ctx context.Context, action *entities.WorkflowAction) error {
	to, ok := action.Config["to"].(string)
	if !ok {
		return errors.NewValidationError("Email recipient not found in config").
			WithField("action.config.to")
	}

	subject, _ := action.Config["subject"].(string)
	body, _ := action.Config["body"].(string)

	s.logger.Info("Email would be sent (email service not configured)", map[string]interface{}{
		"to":      to,
		"subject": subject,
		"body":    body,
	})

	return nil
}

func (s *WorkflowService) executeUpdateRow(ctx context.Context, workflow *entities.Workflow, action *entities.WorkflowAction) error {
	tableName, ok := action.Config["table"].(string)
	if !ok {
		return errors.NewValidationError("Table name not found in config").
			WithField("action.config.table")
	}

	if !validation.IsValidTableName(tableName) {
		return errors.NewValidationError("Invalid table name in workflow action").
			WithField("action.config.table").
			WithDetails(fmt.Sprintf("Table name '%s' is invalid", tableName))
	}

	rowID, ok := action.Config["row_id"].(string)
	if !ok {
		return errors.NewValidationError("Row ID not found in config").
			WithField("action.config.row_id")
	}

	updates, ok := action.Config["updates"].(map[string]interface{})
	if !ok {
		return errors.NewValidationError("Updates not found in config").
			WithField("action.config.updates")
	}

	databases, err := s.dbRepo.GetByProjectID(ctx, workflow.ProjectID)
	if err != nil || len(databases) == 0 {
		return errors.NewNotFoundError("Project database not found")
	}

	database := databases[0]

	return s.dataManipulationSvc.UpdateRow(ctx, database, tableName, rowID, updates)
}

func (s *WorkflowService) executeCreateRow(ctx context.Context, workflow *entities.Workflow, action *entities.WorkflowAction) error {
	tableName, ok := action.Config["table"].(string)
	if !ok {
		return errors.NewValidationError("Table name not found in config").
			WithField("action.config.table")
	}

	if !validation.IsValidTableName(tableName) {
		return errors.NewValidationError("Invalid table name in workflow action").
			WithField("action.config.table").
			WithDetails(fmt.Sprintf("Table name '%s' is invalid", tableName))
	}

	data, ok := action.Config["data"].(map[string]interface{})
	if !ok {
		return errors.NewValidationError("Data not found in config").
			WithField("action.config.data")
	}

	databases, err := s.dbRepo.GetByProjectID(ctx, workflow.ProjectID)
	if err != nil || len(databases) == 0 {
		return errors.NewNotFoundError("Project database not found")
	}

	database := databases[0]

	return s.dataManipulationSvc.CreateRow(ctx, database, tableName, data)
}

func (s *WorkflowService) executeDeleteRow(ctx context.Context, workflow *entities.Workflow, action *entities.WorkflowAction) error {
	tableName, ok := action.Config["table"].(string)
	if !ok {
		return errors.NewValidationError("Table name not found in config").
			WithField("action.config.table")
	}

	if !validation.IsValidTableName(tableName) {
		return errors.NewValidationError("Invalid table name in workflow action").
			WithField("action.config.table").
			WithDetails(fmt.Sprintf("Table name '%s' is invalid", tableName))
	}

	rowID, ok := action.Config["row_id"].(string)
	if !ok {
		return errors.NewValidationError("Row ID not found in config").
			WithField("action.config.row_id")
	}

	databases, err := s.dbRepo.GetByProjectID(ctx, workflow.ProjectID)
	if err != nil || len(databases) == 0 {
		return errors.NewNotFoundError("Project database not found")
	}

	database := databases[0]

	return s.dataManipulationSvc.DeleteRow(ctx, database, tableName, rowID)
}
