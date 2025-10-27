package services

import (
	"context"
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

type UserService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repositories.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(ctx context.Context, req *entities.UserCreateRequest) (*entities.UserResponse, error) {

	if err := validation.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := validation.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	if err := validation.ValidateRequired(req.Name, "name"); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.NewAlreadyExistsError("User with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("Failed to hash password: %w", err))
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
		return nil, errors.NewDatabaseError("Failed to create user", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) Login(ctx context.Context, req *entities.UserLoginRequest) (string, error) {

	if err := validation.ValidateEmail(req.Email); err != nil {
		return "", err
	}

	if err := validation.ValidateRequired(req.Password, "password"); err != nil {
		return "", err
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.NewUnauthorizedError("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.NewUnauthorizedError("Invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.NewInternalError(fmt.Errorf("Failed to generate token: %w", err))
	}

	return tokenString, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("User not found")
	}

	response := user.ToResponse()
	return &response, nil
}
