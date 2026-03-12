package postgres

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	return gormDB, mock
}

func TestNewPostgresAuditLogger(t *testing.T) {
	gormDB, _ := setupTestDB(t)
	logger := NewPostgresAuditLogger(gormDB, "test-service")

	assert.NotNil(t, logger)
	assert.Equal(t, "test-service", logger.serviceName)
	assert.Equal(t, gormDB, logger.db)
}

func TestLog(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	logger := NewPostgresAuditLogger(gormDB, "test-service")

	event := audit.AuditEvent{
		ActorID:      "user-1",
		ActorEmail:   "user@example.com",
		ActorRole:    "admin",
		Action:       "create",
		ResourceType: "user",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "audit"."events"`)).
		WithArgs(
			event.ActorID,
			event.ActorEmail,
			event.ActorRole,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			"test-service",
			event.Action,
			event.ResourceType,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			audit.SeverityInfo,
			audit.CategoryData,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("test-id"))

	err := logger.Log(context.Background(), event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogFromGin(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	logger := NewPostgresAuditLogger(gormDB, "test-service")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest("POST", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("X-Request-ID", "req-123")
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:1234"
	c.Request = req

	c.Set("user_id", "user-1")
	c.Set("email", "test@test.com")
	c.Set("role", "admin")

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "audit"."events"`)).
		WithArgs(
			"user-1",
			"test@test.com",
			"admin",
			"127.0.0.1",
			"test-agent",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			"test-service",
			"test-action",
			"test-resource",
			"test-id",
			sqlmock.AnyArg(),
			"POST",
			"/test",
			"req-123",
			200,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			audit.SeverityInfo,
			audit.CategoryData,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("test-id"))

	err = logger.LogFromGin(c, "test-action", "test-resource", "test-id")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestToDBModel(t *testing.T) {
	event := audit.AuditEvent{
		ActorID:        "user-1",
		ActorEmail:     "user@example.com",
		ActorRole:      "admin",
		ActorIP:        "127.0.0.1",
		ActorUserAgent: "test-agent",
		SchoolID:       "school-1",
		UnitID:         "unit-1",
		ServiceName:    "test-service",
		Action:         "create",
		ResourceType:   "user",
		ResourceID:     "user-123",
		PermissionUsed: "user:create",
		RequestMethod:  "POST",
		RequestPath:    "/users",
		RequestID:      "req-1",
		StatusCode:     201,
		Changes:        map[string]interface{}{"name": "test"},
		Metadata:       map[string]interface{}{"ip": "127.0.0.1"},
		ErrorMessage:   "error",
		Severity:       audit.SeverityWarning,
		Category:       audit.CategoryAuth,
	}

	model := toDBModel(event)

	assert.Equal(t, event.ActorID, model.ActorID)
	assert.Equal(t, event.ActorEmail, model.ActorEmail)
	assert.Equal(t, event.ActorRole, model.ActorRole)
	assert.Equal(t, event.ActorIP, *model.ActorIP)
	assert.Equal(t, event.ActorUserAgent, *model.ActorUserAgent)
	assert.Equal(t, event.SchoolID, *model.SchoolID)
	assert.Equal(t, event.UnitID, *model.UnitID)
	assert.Equal(t, event.ServiceName, model.ServiceName)
	assert.Equal(t, event.Action, model.Action)
	assert.Equal(t, event.ResourceType, model.ResourceType)
	assert.Equal(t, event.ResourceID, *model.ResourceID)
	assert.Equal(t, event.PermissionUsed, *model.PermissionUsed)
	assert.Equal(t, event.RequestMethod, *model.RequestMethod)
	assert.Equal(t, event.RequestPath, *model.RequestPath)
	assert.Equal(t, event.RequestID, *model.RequestID)
	assert.Equal(t, event.StatusCode, *model.StatusCode)
	assert.Equal(t, event.Changes, model.Changes)
	assert.Equal(t, event.Metadata, model.Metadata)
	assert.Equal(t, event.ErrorMessage, *model.ErrorMessage)
	assert.Equal(t, event.Severity, model.Severity)
	assert.Equal(t, event.Category, model.Category)
	assert.Equal(t, "audit.events", model.TableName())
}
