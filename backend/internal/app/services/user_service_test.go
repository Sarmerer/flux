package services

import (
	"context"
	"errors"
	"testing"

	"github.com/flow/internal/app/services/testutil"
	"github.com/flow/internal/app/services/testutil/mocks"
	"github.com/flow/internal/domain/entities"
	apiErrors "github.com/flow/internal/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	req := &entities.UserCreateRequest{
		Email:    "newuser@example.com",
		Password: "SecurePass123!",
		Name:     "New User",
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).
		Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
		return u.Email == req.Email && u.Name == req.Name
	})).Return(nil)

	response, err := service.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Email, response.Email)
	assert.Equal(t, req.Name, response.Name)
	assert.NotEqual(t, uuid.Nil, response.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	existingUser := testutil.NewTestUser()
	req := &entities.UserCreateRequest{
		Email:    existingUser.Email,
		Password: "SecurePass123!",
		Name:     "New User",
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).
		Return(existingUser, nil)

	response, err := service.Register(context.Background(), req)

	assert.Nil(t, response)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeAlreadyExists, apiErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_InvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	testCases := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"invalid format", "notanemail"},
		{"missing domain", "test@"},
		{"missing username", "@example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &entities.UserCreateRequest{
				Email:    tc.email,
				Password: "SecurePass123!",
				Name:     "Test User",
			}

			response, err := service.Register(context.Background(), req)

			assert.Nil(t, response)
			require.Error(t, err)
			apiErr, ok := err.(*apiErrors.APIError)
			require.True(t, ok)
			assert.Equal(t, apiErrors.ErrCodeValidation, apiErr.Code)
		})
	}
}

func TestUserService_Register_WeakPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	testCases := []struct {
		name     string
		password string
	}{
		{"empty password", ""},
		{"too short", "short"},
		{"only letters", "onlyletters"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &entities.UserCreateRequest{
				Email:    "test@example.com",
				Password: tc.password,
				Name:     "Test User",
			}

			response, err := service.Register(context.Background(), req)

			assert.Nil(t, response)
			require.Error(t, err)
			apiErr, ok := err.(*apiErrors.APIError)
			require.True(t, ok)
			assert.Equal(t, apiErrors.ErrCodeValidation, apiErr.Code)
		})
	}
}

func TestUserService_Register_EmptyName(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	req := &entities.UserCreateRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
		Name:     "",
	}

	response, err := service.Register(context.Background(), req)

	assert.Nil(t, response)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeValidation, apiErr.Code)
}

func TestUserService_Register_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	req := &entities.UserCreateRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
		Name:     "Test User",
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).
		Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("database error"))

	response, err := service.Register(context.Background(), req)

	assert.Nil(t, response)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeDatabaseError, apiErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_PasswordIsHashed(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	req := &entities.UserCreateRequest{
		Email:    "test@example.com",
		Password: "PlainTextPassword123!",
		Name:     "Test User",
	}

	var capturedUser *entities.User
	mockRepo.On("GetByEmail", mock.Anything, req.Email).
		Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*entities.User)
		}).
		Return(nil)

	_, err := service.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotEqual(t, req.Password, capturedUser.Password)
	err = bcrypt.CompareHashAndPassword([]byte(capturedUser.Password), []byte(req.Password))
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	password := "SecurePass123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := testutil.NewTestUser()
	user.Password = string(hashedPassword)

	req := &entities.UserLoginRequest{
		Email:    user.Email,
		Password: password,
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).
		Return(user, nil)

	token, err := service.Login(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := testutil.NewTestUser()
	user.Password = string(hashedPassword)

	req := &entities.UserLoginRequest{
		Email:    user.Email,
		Password: "wrongpassword",
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).
		Return(user, nil)

	token, err := service.Login(context.Background(), req)

	assert.Empty(t, token)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeUnauthorized, apiErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	req := &entities.UserLoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password",
	}

	mockRepo.On("GetByEmail", mock.Anything, req.Email).
		Return(nil, errors.New("not found"))

	token, err := service.Login(context.Background(), req)

	assert.Empty(t, token)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeUnauthorized, apiErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	req := &entities.UserLoginRequest{
		Email:    "invalidemail",
		Password: "password",
	}

	token, err := service.Login(context.Background(), req)

	assert.Empty(t, token)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeValidation, apiErr.Code)
}

func TestUserService_Login_EmptyPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	req := &entities.UserLoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	token, err := service.Login(context.Background(), req)

	assert.Empty(t, token)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeValidation, apiErr.Code)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	user := testutil.NewTestUser()

	mockRepo.On("GetByID", mock.Anything, user.ID).
		Return(user, nil)

	response, err := service.GetUserByID(context.Background(), user.ID)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo, "test-secret")

	userID := uuid.New()

	mockRepo.On("GetByID", mock.Anything, userID).
		Return(nil, errors.New("not found"))

	response, err := service.GetUserByID(context.Background(), userID)

	assert.Nil(t, response)
	require.Error(t, err)
	apiErr, ok := err.(*apiErrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apiErrors.ErrCodeNotFound, apiErr.Code)
	mockRepo.AssertExpectations(t)
}
