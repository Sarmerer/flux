package testutil

import (
	"time"

	"github.com/flow/internal/domain/entities"
	"github.com/google/uuid"
)

func NewTestUser() *entities.User {
	return &entities.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTestUserWithEmail(email string) *entities.User {
	user := NewTestUser()
	user.Email = email
	return user
}

func NewTestProject(ownerID uuid.UUID) *entities.Project {
	return &entities.Project{
		ID:          uuid.New(),
		Name:        "Test Project",
		Description: "Test project description",
		OwnerID:     ownerID,
		APIKey:      "test-api-key-12345",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewTestDatabase(projectID uuid.UUID) *entities.Database {
	return &entities.Database{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      "test_db",
		Host:      "localhost",
		Port:      5432,
		Username:  "testuser",
		Password:  "testpass",
		Database:  "test_project_db",
		SSLMode:   "disable",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTestTable(projectID uuid.UUID) *entities.Table {
	return &entities.Table{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        "test_table",
		Description: "Test table description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewTestWorkflow(projectID uuid.UUID) *entities.Workflow {
	return &entities.Workflow{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        "Test Workflow",
		Description: "Test workflow description",
		Trigger: entities.WorkflowTrigger{
			Type: "manual",
		},
		Actions: []entities.WorkflowAction{
			{
				ID:   "action-1",
				Type: "send_webhook",
				Config: map[string]interface{}{
					"url": "https://example.com/webhook",
				},
			},
		},
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTestProjectMember(projectID, userID uuid.UUID, role entities.Role) *entities.ProjectMember {
	return &entities.ProjectMember{
		ID:        uuid.New(),
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		Permissions: entities.Permissions{
			entities.PermProjectView,
			entities.PermTableEdit,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTestTableSchemaInfo() *entities.TableSchemaInfo {
	return &entities.TableSchemaInfo{
		Columns: []entities.ColumnInfo{
			{
				Name:            "id",
				DataType:        "uuid",
				IsNullable:      false,
				OrdinalPosition: 1,
			},
			{
				Name:            "name",
				DataType:        "varchar",
				IsNullable:      false,
				OrdinalPosition: 2,
			},
			{
				Name:            "email",
				DataType:        "varchar",
				IsNullable:      true,
				DefaultValue:    nil,
				OrdinalPosition: 3,
			},
		},
		PrimaryKeys: []string{"id"},
		ForeignKeys: []entities.ForeignKeyInfo{},
		Indexes: []entities.IndexInfo{
			{
				Name:        "test_table_pkey",
				ColumnNames: []string{"id"},
				IsUnique:    true,
				IsPrimary:   true,
				IndexType:   "btree",
			},
		},
	}
}
