package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flow/internal/domain/entities"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, req *entities.UserCreateRequest) (*entities.UserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*entities.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *entities.UserLoginRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.UserResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.UserResponse), args.Error(1)
}

type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) CreateProject(ctx context.Context, req *entities.ProjectCreateRequest, userID uuid.UUID) (*entities.ProjectResponse, error) {
	args := m.Called(ctx, req, userID)
	return args.Get(0).(*entities.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*entities.ProjectResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) GetProjectsByOwnerID(ctx context.Context, userID uuid.UUID) ([]*entities.ProjectResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) UpdateProject(ctx context.Context, id uuid.UUID, req *entities.ProjectCreateRequest) (*entities.ProjectResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*entities.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockTableService struct {
	mock.Mock
}

func (m *MockTableService) CreateTable(ctx context.Context, req *entities.TableCreateRequest, projectID uuid.UUID) (*entities.TableResponse, error) {
	args := m.Called(ctx, req, projectID)
	return args.Get(0).(*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) GetTableByID(ctx context.Context, id uuid.UUID) (*entities.TableResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) GetTablesByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entities.TableResponse, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) UpdateTable(ctx context.Context, id uuid.UUID, req *entities.TableUpdateRequest) (*entities.TableResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*entities.TableResponse), args.Error(1)
}

func (m *MockTableService) DeleteTable(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_Register(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewUserHandler(mockUserService)

	userReq := &entities.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	expectedUser := &entities.UserResponse{
		ID:    uuid.New(),
		Email: userReq.Email,
		Name:  userReq.Name,
	}

	mockUserService.On("Register", mock.Anything, userReq).Return(expectedUser, nil)

	jsonData, _ := json.Marshal(userReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response entities.User
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, response.Email)
	assert.Equal(t, expectedUser.Name, response.Name)

	mockUserService.AssertExpectations(t)
}

func TestUserHandler_Login(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewUserHandler(mockUserService)

	loginReq := &entities.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedToken := "jwt-token-here"

	mockUserService.On("Login", mock.Anything, loginReq).Return(expectedToken, nil)

	jsonData, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, response["token"])

	mockUserService.AssertExpectations(t)
}

func TestProjectHandler_CreateProject(t *testing.T) {
	mockProjectService := new(MockProjectService)
	handler := NewProjectHandler(mockProjectService)

	projectReq := &entities.ProjectCreateRequest{
		Name:        "Test Project",
		Description: "A test project",
	}

	expectedProject := &entities.ProjectResponse{
		ID:          uuid.New(),
		Name:        projectReq.Name,
		Description: projectReq.Description,
		OwnerID:     uuid.New(),
	}

	mockProjectService.On("CreateProject", mock.Anything, projectReq, mock.AnythingOfType("uuid.UUID")).Return(expectedProject, nil)

	jsonData, _ := json.Marshal(projectReq)
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), "user_id", uuid.New())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.CreateProject(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response entities.Project
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedProject.Name, response.Name)
	assert.Equal(t, expectedProject.Description, response.Description)

	mockProjectService.AssertExpectations(t)
}

func TestTableHandler_CreateTable(t *testing.T) {
	mockTableService := new(MockTableService)
	handler := NewTableHandler(mockTableService)

	tableReq := &entities.TableCreateRequest{
		Name: "test_table",
		Schema: map[string]interface{}{
			"columns": []map[string]interface{}{
				{"name": "id", "type": "uuid"},
			},
		},
	}

	projectID := uuid.New()
	expectedTable := &entities.TableResponse{
		ID:        uuid.New(),
		Name:      tableReq.Name,
		Schema:    `{"columns": [{"name": "id", "type": "uuid"}]}`,
		ProjectID: projectID,
	}

	mockTableService.On("CreateTable", mock.Anything, mock.AnythingOfType("*entities.TableCreateRequest"), projectID).Return(expectedTable, nil)

	jsonData, _ := json.Marshal(tableReq)

	r := chi.NewRouter()
	r.Post("/api/v1/projects/{projectId}/tables", handler.CreateTable)

	req, _ := http.NewRequest("POST", "/api/v1/projects/"+projectID.String()+"/tables", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response entities.TableResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTable.Name, response.Name)
	assert.Equal(t, expectedTable.Schema, response.Schema)

	mockTableService.AssertExpectations(t)
}

func TestUserHandler_Register_InvalidJSON(t *testing.T) {
	mockUserService := new(MockUserService)
	handler := NewUserHandler(mockUserService)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Register(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestProjectHandler_CreateProject_Unauthorized(t *testing.T) {
	mockProjectService := new(MockProjectService)
	handler := NewProjectHandler(mockProjectService)

	projectReq := &entities.ProjectCreateRequest{
		Name: "Test Project",
	}
	jsonData, _ := json.Marshal(projectReq)
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateProject(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
