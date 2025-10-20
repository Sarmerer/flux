package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
	// Validation errors
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField  ErrorCode = "MISSING_FIELD"
	ErrCodeInvalidFormat ErrorCode = "INVALID_FORMAT"

	// Authentication errors
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeInvalidToken ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"

	// Resource errors
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict      ErrorCode = "CONFLICT"

	// Database errors
	ErrCodeDatabaseError    ErrorCode = "DATABASE_ERROR"
	ErrCodeConnectionError  ErrorCode = "CONNECTION_ERROR"
	ErrCodeQueryError       ErrorCode = "QUERY_ERROR"
	ErrCodeTransactionError ErrorCode = "TRANSACTION_ERROR"

	// Business logic errors
	ErrCodeBusinessRule    ErrorCode = "BUSINESS_RULE_VIOLATION"
	ErrCodeInvalidState    ErrorCode = "INVALID_STATE"
	ErrCodeOperationFailed ErrorCode = "OPERATION_FAILED"

	// External service errors
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeTimeout         ErrorCode = "TIMEOUT"

	// Internal errors
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"
)

// APIError represents a structured API error response
type APIError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Field     string                 `json:"field,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError
func NewAPIError(code ErrorCode, message string) *APIError {
	return &APIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// WithDetails adds details to the error
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// WithField adds a field name to the error
func (e *APIError) WithField(field string) *APIError {
	e.Field = field
	return e
}

// WithRequestID adds a request ID to the error
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}

// WithMetadata adds metadata to the error
func (e *APIError) WithMetadata(metadata map[string]interface{}) *APIError {
	e.Metadata = metadata
	return e
}

// HTTPStatus returns the appropriate HTTP status code for the error
func (e *APIError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeValidation, ErrCodeInvalidInput, ErrCodeMissingField, ErrCodeInvalidFormat:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeTokenExpired:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists, ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeDatabaseError, ErrCodeConnectionError, ErrCodeQueryError, ErrCodeTransactionError:
		return http.StatusInternalServerError
	case ErrCodeBusinessRule, ErrCodeInvalidState, ErrCodeOperationFailed:
		return http.StatusUnprocessableEntity
	case ErrCodeExternalService, ErrCodeTimeout:
		return http.StatusBadGateway
	case ErrCodeNotImplemented:
		return http.StatusNotImplemented
	case ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// WriteError writes an error response to the HTTP response writer
func WriteError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus())

	json.NewEncoder(w).Encode(err)
}

// WriteErrorResponse is a convenience function for writing error responses
func WriteErrorResponse(w http.ResponseWriter, code ErrorCode, message string) {
	err := NewAPIError(code, message)
	WriteError(w, err)
}

// Common error constructors
func NewValidationError(message string) *APIError {
	return NewAPIError(ErrCodeValidation, message)
}

func NewNotFoundError(resource string) *APIError {
	return NewAPIError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

func NewUnauthorizedError() *APIError {
	return NewAPIError(ErrCodeUnauthorized, "Unauthorized access")
}

func NewForbiddenError() *APIError {
	return NewAPIError(ErrCodeForbidden, "Access forbidden")
}

func NewDatabaseError(err error) *APIError {
	return NewAPIError(ErrCodeDatabaseError, "Database operation failed").WithDetails(err.Error())
}

func NewInternalError(err error) *APIError {
	return NewAPIError(ErrCodeInternal, "Internal server error").WithDetails(err.Error())
}

func NewBusinessRuleError(message string) *APIError {
	return NewAPIError(ErrCodeBusinessRule, message)
}

func NewConflictError(message string) *APIError {
	return NewAPIError(ErrCodeConflict, message)
}

func NewAlreadyExistsError(resource string) *APIError {
	return NewAPIError(ErrCodeAlreadyExists, fmt.Sprintf("%s already exists", resource))
}

// ErrorHandler is a middleware for handling panics and converting them to API errors
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				apiErr := NewInternalError(fmt.Errorf("panic: %v", err))
				WriteError(w, apiErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
