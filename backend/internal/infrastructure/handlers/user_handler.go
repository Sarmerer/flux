package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/infrastructure/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entities.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.NewValidationError("Invalid request body").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		errors.WriteError(w, errors.NewValidationError("Email is required").WithField("email"))
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		errors.WriteError(w, errors.NewValidationError("Password is required").WithField("password"))
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		errors.WriteError(w, errors.NewValidationError("Name is required").WithField("name"))
		return
	}

	user, err := h.userService.Register(r.Context(), &req)
	if err != nil {

		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewInternalError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entities.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.NewValidationError("Invalid request body").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		errors.WriteError(w, errors.NewValidationError("Email is required").WithField("email"))
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		errors.WriteError(w, errors.NewValidationError("Password is required").WithField("password"))
		return
	}

	token, err := h.userService.Login(r.Context(), &req)
	if err != nil {

		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewUnauthorizedError().WithDetails(err.Error()))
		}
		return
	}

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		apiErr := errors.NewValidationError("Invalid user ID").WithField("id").WithDetails(err.Error())
		errors.WriteError(w, apiErr)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewNotFoundError("User"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError().WithDetails("User ID not found in context"))
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			errors.WriteError(w, apiErr)
		} else {
			errors.WriteError(w, errors.NewNotFoundError("User"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
