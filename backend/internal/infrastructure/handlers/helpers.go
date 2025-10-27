package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flow/internal/validation"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func parseUUIDParam(r *http.Request, paramName string) (uuid.UUID, error) {
	paramStr := chi.URLParam(r, paramName)
	return uuid.Parse(paramStr)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(validation.Errors); ok {
		return true
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "validation failed") ||
		strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "required") ||
		strings.Contains(errMsg, "duplicate")
}
