package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/flow/internal/errors"
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

func writeError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*errors.APIError); ok {
		errors.WriteError(w, apiErr)
	} else {
		errors.WriteError(w, errors.NewInternalError(err))
	}
}
