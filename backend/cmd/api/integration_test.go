package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flow/internal/app/services"
	"github.com/flow/internal/domain/entities"
	postgres "github.com/flow/internal/infrastructure/database"
	"github.com/flow/internal/infrastructure/handlers"
	postgresRepo "github.com/flow/internal/infrastructure/repositories/postgres"
	"github.com/flow/internal/infrastructure/routes"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IntegrationTestSuite provides a test suite for integration testing
type IntegrationTestSuite struct {
	router http.Handler
	db     *pgxpool.Pool
	userID uuid.UUID
	token  string
}

// SetupTestSuite initializes the test suite with a test database
func SetupTestSuite(t *testing.T) *IntegrationTestSuite {
	// Use test database configuration
	dbConfig := postgres.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "flow_test", // Use test database
		SSLMode:  "disable",
	}

	ctx := context.Background()
	db, err := postgres.NewConnection(ctx, dbConfig)
	require.NoError(t, err, "Failed to connect to test database")

	// Run migrations
	err = postgres.Migrate(ctx, db)
	require.NoError(t, err, "Failed to run migrations")

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(db)
	projectRepo := postgresRepo.NewProjectRepository(db)
	databaseRepo := postgresRepo.NewDatabaseRepository(db)
	tableRepo := postgresRepo.NewTableRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, "test-jwt-secret")
	projectService := services.NewProjectService(projectRepo, databaseRepo)
	tableService := services.NewTableService(tableRepo)

	// Initialize mutation services
	dbMutationService := services.NewDatabaseMutationService(databaseRepo, projectRepo, db)
	tableSchemaMutationService := services.NewTableSchemaMutationService(tableRepo, databaseRepo, projectRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	tableHandler := handlers.NewTableHandler(tableService)
	dbMutationHandler := handlers.NewDatabaseMutationHandler(dbMutationService)
	tableSchemaMutationHandler := handlers.NewTableSchemaMutationHandler(tableSchemaMutationService)

	// Setup routes
	router := routes.SetupRoutes(userHandler, projectHandler, tableHandler, dbMutationHandler, tableSchemaMutationHandler)

	return &IntegrationTestSuite{
		router: router,
		db:     db,
	}
}

// CleanupTestSuite cleans up test data
func (suite *IntegrationTestSuite) CleanupTestSuite(t *testing.T) {
	if suite.db != nil {
		suite.db.Close()
	}
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.CleanupTestSuite(t)

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

// TestUserRegistrationAndLogin tests the complete user flow
func TestUserRegistrationAndLogin(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.CleanupTestSuite(t)

	// Test user registration
	userData := entities.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	jsonData, err := json.Marshal(userData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var userResponse entities.User
	err = json.Unmarshal(rr.Body.Bytes(), &userResponse)
	require.NoError(t, err)
	assert.Equal(t, userData.Email, userResponse.Email)
	assert.Equal(t, userData.Name, userResponse.Name)

	// Test user login
	loginData := entities.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var loginResponse map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResponse["token"])

	// Store token for further tests
	suite.token = loginResponse["token"]
	suite.userID = userResponse.ID
}

// TestProjectManagement tests project CRUD operations
func TestProjectManagement(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.CleanupTestSuite(t)

	// First register and login a user
	suite.setupUser(t)

	// Test project creation
	projectData := entities.ProjectCreateRequest{
		Name:        "Test Project",
		Description: "A test project",
	}

	jsonData, err := json.Marshal(projectData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var projectResponse entities.Project
	err = json.Unmarshal(rr.Body.Bytes(), &projectResponse)
	require.NoError(t, err)
	assert.Equal(t, projectData.Name, projectResponse.Name)
	assert.Equal(t, projectData.Description, projectResponse.Description)

	// Test getting projects
	req, err = http.NewRequest("GET", "/api/v1/projects", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	rr = httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var projectsResponse []entities.Project
	err = json.Unmarshal(rr.Body.Bytes(), &projectsResponse)
	require.NoError(t, err)
	assert.Len(t, projectsResponse, 1)
	assert.Equal(t, projectData.Name, projectsResponse[0].Name)
}

// TestTableManagement tests table CRUD operations
func TestTableManagement(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.CleanupTestSuite(t)

	// Setup user and project
	projectID := suite.setupUserAndProject(t)

	// Test table creation
	tableData := entities.TableCreateRequest{
		Name: "test_table",
		Schema: map[string]interface{}{
			"columns": []map[string]interface{}{
				{"name": "id", "type": "uuid"},
				{"name": "name", "type": "varchar"},
			},
		},
	}

	jsonData, err := json.Marshal(tableData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/projects/"+projectID.String()+"/tables", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var tableResponse entities.Table
	err = json.Unmarshal(rr.Body.Bytes(), &tableResponse)
	require.NoError(t, err)
	assert.Equal(t, tableData.Name, tableResponse.Name)
	assert.Equal(t, tableData.Schema, tableResponse.Schema)

	// Test getting tables
	req, err = http.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/tables", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	rr = httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var tablesResponse []entities.Table
	err = json.Unmarshal(rr.Body.Bytes(), &tablesResponse)
	require.NoError(t, err)
	assert.Len(t, tablesResponse, 1)
	assert.Equal(t, tableData.Name, tablesResponse[0].Name)
}

// Helper methods
func (suite *IntegrationTestSuite) setupUser(t *testing.T) {
	userData := entities.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	var userResponse entities.User
	json.Unmarshal(rr.Body.Bytes(), &userResponse)
	suite.userID = userResponse.ID

	// Login to get token
	loginData := entities.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonData, _ = json.Marshal(loginData)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	var loginResponse map[string]string
	json.Unmarshal(rr.Body.Bytes(), &loginResponse)
	suite.token = loginResponse["token"]
}

func (suite *IntegrationTestSuite) setupUserAndProject(t *testing.T) uuid.UUID {
	suite.setupUser(t)

	projectData := entities.ProjectCreateRequest{
		Name:        "Test Project",
		Description: "A test project",
	}

	jsonData, _ := json.Marshal(projectData)
	req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	var projectResponse entities.Project
	json.Unmarshal(rr.Body.Bytes(), &projectResponse)
	return projectResponse.ID
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.CleanupTestSuite(t)

	// Test invalid JSON
	req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test unauthorized access
	req, err = http.NewRequest("GET", "/api/v1/projects", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// Benchmark tests
func BenchmarkHealthEndpoint(b *testing.B) {
	suite := SetupTestSuite(&testing.T{})
	defer suite.CleanupTestSuite(&testing.T{})

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.router.ServeHTTP(rr, req)
	}
}
