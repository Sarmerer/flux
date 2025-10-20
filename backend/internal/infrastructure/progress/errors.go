package progress

import "errors"

var (
	ErrOperationNotFound  = errors.New("operation not found")
	ErrInvalidProgress    = errors.New("invalid progress value (must be between 0.0 and 1.0)")
	ErrOperationCompleted = errors.New("operation already completed")
	ErrOperationCancelled = errors.New("operation was cancelled")
)
