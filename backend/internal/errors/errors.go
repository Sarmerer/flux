package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/flow/internal/infrastructure/logging"
)

type ErrorCode string

const (
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField  ErrorCode = "MISSING_FIELD"
	ErrCodeInvalidFormat ErrorCode = "INVALID_FORMAT"

	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeInvalidToken ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"

	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict      ErrorCode = "CONFLICT"

	ErrCodeDatabaseError    ErrorCode = "DATABASE_ERROR"
	ErrCodeConnectionError  ErrorCode = "CONNECTION_ERROR"
	ErrCodeQueryError       ErrorCode = "QUERY_ERROR"
	ErrCodeTransactionError ErrorCode = "TRANSACTION_ERROR"

	ErrCodeBusinessRule    ErrorCode = "BUSINESS_RULE_VIOLATION"
	ErrCodeInvalidState    ErrorCode = "INVALID_STATE"
	ErrCodeOperationFailed ErrorCode = "OPERATION_FAILED"

	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeTimeout         ErrorCode = "TIMEOUT"

	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"
)

type APIError struct {
	Code          ErrorCode              `json:"code"`
	Message       string                 `json:"message"`
	Details       string                 `json:"details,omitempty"`
	Field         string                 `json:"field,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	StackTrace    string                 `json:"-"`
	OriginalError error                  `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code ErrorCode, message string) *APIError {
	return &APIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

func (e *APIError) WithField(field string) *APIError {
	e.Field = field
	return e
}

func (e *APIError) WithMetadata(metadata map[string]interface{}) *APIError {
	e.Metadata = metadata
	return e
}

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

func WriteError(w http.ResponseWriter, err *APIError, logger ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus())

	var loggerInstance interface{}
	if len(logger) > 0 {
		loggerInstance = logger[0]
	}

	if err.HTTPStatus() >= 500 && loggerInstance != nil {
		if ctxLogger, ok := loggerInstance.(interface {
			ErrorWithStack(msg string, err error, stackTrace string, fields ...map[string]interface{})
		}); ok {
			fields := map[string]interface{}{
				"error_code": err.Code,
				"details":    err.Details,
			}

			if err.Field != "" {
				fields["field"] = err.Field
			}

			if len(err.Metadata) > 0 {
				fields["metadata"] = err.Metadata
			}

			ctxLogger.ErrorWithStack(err.Message, err.OriginalError, err.StackTrace, fields)
		}
	}

	json.NewEncoder(w).Encode(err)
}

func NewValidationError(message string) *APIError {
	return NewAPIError(ErrCodeValidation, message)
}

func NewNotFoundError(message string) *APIError {
	return NewAPIError(ErrCodeNotFound, message)
}

func NewUnauthorizedError(message string) *APIError {
	return NewAPIError(ErrCodeUnauthorized, message)
}

func NewForbiddenError(message string) *APIError {
	return NewAPIError(ErrCodeForbidden, message)
}

func NewDatabaseError(message string, err error) *APIError {
	apiErr := NewAPIError(ErrCodeDatabaseError, message)
	if err != nil {
		apiErr.Details = err.Error()
		apiErr.OriginalError = err
	}
	return apiErr
}

func NewInternalError(err error) *APIError {
	message := "Internal server error"
	if err != nil {
		apiErr := NewAPIError(ErrCodeInternal, message).WithDetails(err.Error())
		apiErr.OriginalError = err
		apiErr.StackTrace = string(debug.Stack())
		return apiErr
	}
	return NewAPIError(ErrCodeInternal, message)
}

func NewBusinessRuleError(message string) *APIError {
	return NewAPIError(ErrCodeBusinessRule, message)
}

func NewConflictError(message string) *APIError {
	return NewAPIError(ErrCodeConflict, message)
}

func NewAlreadyExistsError(message string) *APIError {
	return NewAPIError(ErrCodeAlreadyExists, message)
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				apiErr := NewInternalError(fmt.Errorf("panic: %v", err))

				logger := logging.GetLogger(r.Context())
				if logger != nil {
					fields := map[string]interface{}{
						"panic_value": err,
						"method":      r.Method,
						"path":        r.URL.Path,
						"remote_addr": r.RemoteAddr,
					}
					logger.ErrorWithStack("Panic recovered", apiErr.OriginalError, apiErr.StackTrace, fields)
				}

				WriteError(w, apiErr, logger)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
