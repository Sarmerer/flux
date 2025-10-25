package progress

import (
	"context"
	"sync"
	"time"

	"github.com/flow/internal/infrastructure/realtime"
	"github.com/google/uuid"
)

type OperationStatus string

const (
	StatusPending   OperationStatus = "pending"
	StatusRunning   OperationStatus = "running"
	StatusCompleted OperationStatus = "completed"
	StatusFailed    OperationStatus = "failed"
	StatusCancelled OperationStatus = "cancelled"
)

type Operation struct {
	ID          string                 `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	Type        string                 `json:"type"`
	Status      OperationStatus        `json:"status"`
	Progress    float64                `json:"progress"`
	CurrentStep string                 `json:"current_step"`
	Message     string                 `json:"message"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Result      interface{}            `json:"result,omitempty"`
	Steps       []Step                 `json:"steps"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type Step struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      OperationStatus `json:"status"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     *time.Time      `json:"end_time,omitempty"`
	Progress    float64         `json:"progress"`
	Error       string          `json:"error,omitempty"`
}

type Tracker struct {
	operations map[string]*Operation
	mutex      sync.RWMutex
	hub        *realtime.Hub
}

func NewTracker(hub *realtime.Hub) *Tracker {
	return &Tracker{
		operations: make(map[string]*Operation),
		hub:        hub,
	}
}

func (t *Tracker) StartOperation(userID uuid.UUID, operationType, message string, steps []string) *Operation {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	operationID := uuid.New().String()

	operationSteps := make([]Step, len(steps))
	for i, stepName := range steps {
		operationSteps[i] = Step{
			Name:        stepName,
			Description: stepName,
			Status:      StatusPending,
			Progress:    0.0,
		}
	}

	operation := &Operation{
		ID:          operationID,
		UserID:      userID,
		Type:        operationType,
		Status:      StatusPending,
		Progress:    0.0,
		CurrentStep: "",
		Message:     message,
		StartTime:   time.Now(),
		Steps:       operationSteps,
		Metadata:    make(map[string]interface{}),
	}

	t.operations[operationID] = operation

	realtime.SendNotificationToUser(t.hub, userID, "Operation started", map[string]interface{}{
		"operation_id": operationID,
		"type":         operationType,
		"message":      message,
	})

	return operation
}

func (t *Tracker) UpdateProgress(operationID, stepName, message string, progress float64) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	operation, exists := t.operations[operationID]
	if !exists {
		return ErrOperationNotFound
	}

	operation.Progress = progress
	operation.CurrentStep = stepName
	operation.Message = message
	operation.Status = StatusRunning

	for i, step := range operation.Steps {
		if step.Name == stepName {
			operation.Steps[i].Status = StatusRunning
			operation.Steps[i].Progress = progress
			if operation.Steps[i].StartTime.IsZero() {
				operation.Steps[i].StartTime = time.Now()
			}
			break
		}
	}

	realtime.SendProgressToUser(t.hub, operation.UserID, operationID, stepName, progress, message)

	return nil
}

func (t *Tracker) CompleteStep(operationID, stepName string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	operation, exists := t.operations[operationID]
	if !exists {
		return ErrOperationNotFound
	}

	for i, step := range operation.Steps {
		if step.Name == stepName {
			operation.Steps[i].Status = StatusCompleted
			operation.Steps[i].Progress = 1.0
			now := time.Now()
			operation.Steps[i].EndTime = &now
			break
		}
	}

	return nil
}

func (t *Tracker) CompleteOperation(operationID string, result interface{}) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	operation, exists := t.operations[operationID]
	if !exists {
		return ErrOperationNotFound
	}

	operation.Status = StatusCompleted
	operation.Progress = 1.0
	operation.Result = result
	now := time.Now()
	operation.EndTime = &now

	realtime.SendSuccessToUser(t.hub, operation.UserID, operationID, "Operation completed successfully", result)

	return nil
}

func (t *Tracker) FailOperation(operationID string, err error) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	operation, exists := t.operations[operationID]
	if !exists {
		return ErrOperationNotFound
	}

	operation.Status = StatusFailed
	operation.Error = err.Error()
	now := time.Now()
	operation.EndTime = &now

	realtime.SendErrorToUser(t.hub, operation.UserID, operationID, "Operation failed", err)

	return nil
}

func (t *Tracker) CancelOperation(operationID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	operation, exists := t.operations[operationID]
	if !exists {
		return ErrOperationNotFound
	}

	operation.Status = StatusCancelled
	now := time.Now()
	operation.EndTime = &now

	realtime.SendNotificationToUser(t.hub, operation.UserID, "Operation cancelled", map[string]interface{}{
		"operation_id": operationID,
		"status":       StatusCancelled,
	})

	return nil
}

func (t *Tracker) GetOperation(operationID string) (*Operation, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	operation, exists := t.operations[operationID]
	if !exists {
		return nil, ErrOperationNotFound
	}

	return operation, nil
}

func (t *Tracker) GetUserOperations(userID uuid.UUID) []*Operation {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	var userOperations []*Operation
	for _, operation := range t.operations {
		if operation.UserID == userID {
			userOperations = append(userOperations, operation)
		}
	}

	return userOperations
}

func (t *Tracker) CleanupCompletedOperations(olderThan time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	cutoff := time.Now().Add(-olderThan)
	for id, operation := range t.operations {
		if (operation.Status == StatusCompleted || operation.Status == StatusFailed || operation.Status == StatusCancelled) &&
			operation.EndTime != nil && operation.EndTime.Before(cutoff) {
			delete(t.operations, id)
		}
	}
}

func (t *Tracker) StartCleanupRoutine(ctx context.Context, interval, olderThan time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t.CleanupCompletedOperations(olderThan)
			}
		}
	}()
}

type OperationContext struct {
	OperationID string
	UserID      uuid.UUID
	Tracker     *Tracker
	ctx         context.Context
	cancel      context.CancelFunc
}

func (t *Tracker) NewOperationContext(userID uuid.UUID, operationType, message string, steps []string) *OperationContext {
	operation := t.StartOperation(userID, operationType, message, steps)
	ctx, cancel := context.WithCancel(context.Background())

	return &OperationContext{
		OperationID: operation.ID,
		UserID:      userID,
		Tracker:     t,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (oc *OperationContext) UpdateProgress(stepName, message string, progress float64) error {
	return oc.Tracker.UpdateProgress(oc.OperationID, stepName, message, progress)
}

func (oc *OperationContext) CompleteStep(stepName string) error {
	return oc.Tracker.CompleteStep(oc.OperationID, stepName)
}

func (oc *OperationContext) Complete(result interface{}) error {
	oc.cancel()
	return oc.Tracker.CompleteOperation(oc.OperationID, result)
}

func (oc *OperationContext) Fail(err error) error {
	oc.cancel()
	return oc.Tracker.FailOperation(oc.OperationID, err)
}

func (oc *OperationContext) Cancel() error {
	oc.cancel()
	return oc.Tracker.CancelOperation(oc.OperationID)
}

func (oc *OperationContext) Context() context.Context {
	return oc.ctx
}

func (oc *OperationContext) Done() <-chan struct{} {
	return oc.ctx.Done()
}
