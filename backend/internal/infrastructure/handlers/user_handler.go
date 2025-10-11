package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entities.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entities.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile handles getting user profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
