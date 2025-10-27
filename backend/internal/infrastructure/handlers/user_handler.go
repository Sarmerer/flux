package handlers

import (
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
	if err := decodeJSON(r, &req); err != nil {
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
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entities.UserLoginRequest
	if err := decodeJSON(r, &req); err != nil {
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
		writeError(w, err)
		return
	}

	response := map[string]string{"token": token}
	writeJSON(w, http.StatusOK, response)
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
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
