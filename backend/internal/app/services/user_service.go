package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/flow/internal/domain/repositories"
	"github.com/flow/internal/errors"
	"github.com/flow/internal/validation"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

// NewUserService creates a new UserService
func NewUserService(userRepo repositories.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user
func (s *UserService) Register(ctx context.Context, req *entities.UserCreateRequest) (*entities.UserResponse, error) {
	// Validate email
	if err := validation.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	// Validate password strength
	if err := validation.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Validate name
	if err := validation.ValidateRequired(req.Name, "name"); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.NewAlreadyExistsError("User with this email")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to hash password: %w", err))
	}

	user := &entities.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	response := user.ToResponse()
	return &response, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(ctx context.Context, req *entities.UserLoginRequest) (string, error) {
	// Validate email
	if err := validation.ValidateEmail(req.Email); err != nil {
		return "", err
	}

	// Validate password is provided
	if err := validation.ValidateRequired(req.Password, "password"); err != nil {
		return "", err
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.NewUnauthorizedError().WithDetails("Invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.NewUnauthorizedError().WithDetails("Invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.NewInternalError(fmt.Errorf("failed to generate token: %w", err))
	}

	return tokenString, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("User")
	}

	response := user.ToResponse()
	return &response, nil
}

// generateAPIKey generates a random API key
func generateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
