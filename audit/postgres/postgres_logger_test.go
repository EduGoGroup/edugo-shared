package postgres

import (
	"context"
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
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	return gormDB, mock
}

func TestPostgresAuditLogger_Log(t *testing.T) {
	db, mock := setupMockDB(t)
	loggerInstance := NewPostgresAuditLogger(db, "test-service")

	event := audit.AuditEvent{
		ActorID:      "user-123",
		ActorEmail:   "test@example.com",
		ActorRole:    "admin",
		Action:       "CREATE",
		ResourceType: "User",
		Severity:     audit.SeverityInfo,
		Category:     audit.CategoryData,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "audit"."events"`)).
		WithArgs(event.ActorID, event.ActorEmail, event.ActorRole, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "test-service", event.Action, event.ResourceType, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), event.Severity, event.Category).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))
	mock.ExpectCommit()

	err := loggerInstance.Log(context.Background(), event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresAuditLogger_LogFromGin(t *testing.T) {
	db, mock := setupMockDB(t)
	loggerInstance := NewPostgresAuditLogger(db, "test-service")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/users", nil)
	c.Request.Header.Set("X-Request-ID", "req-123")
	c.Request.Header.Set("User-Agent", "test-agent")
	c.Set("user_id", "user-123")
	c.Set("email", "test@example.com")
	c.Set("role", "admin")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "audit"."events"`)).
		WithArgs("user-123", "test@example.com", "admin", sqlmock.AnyArg(), "test-agent", sqlmock.AnyArg(), sqlmock.AnyArg(), "test-service", "CREATE", "User", "user-456", sqlmock.AnyArg(), "POST", "/users", "req-123", 200, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), audit.SeverityInfo, audit.CategoryData).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))
	mock.ExpectCommit()

	err := loggerInstance.LogFromGin(c, "CREATE", "User", "user-456")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestToDBModel(t *testing.T) {
	event := audit.AuditEvent{
		ActorID:      "user-123",
		ActorEmail:   "test@example.com",
		ActorRole:    "admin",
		Action:       "CREATE",
		ResourceType: "User",
		ResourceID:   "user-456",
		ServiceName:  "test-service",
		Severity:     audit.SeverityInfo,
		Category:     audit.CategoryData,
		ActorIP:      "127.0.0.1",
	}

	model := toDBModel(event)

	assert.Equal(t, "user-123", model.ActorID)
	assert.Equal(t, "test@example.com", model.ActorEmail)
	assert.Equal(t, "admin", model.ActorRole)
	assert.Equal(t, "CREATE", model.Action)
	assert.Equal(t, "User", model.ResourceType)
	assert.NotNil(t, model.ResourceID)
	assert.Equal(t, "user-456", *model.ResourceID)
	assert.Equal(t, "test-service", model.ServiceName)
	assert.Equal(t, audit.SeverityInfo, model.Severity)
	assert.Equal(t, audit.CategoryData, model.Category)
	assert.NotNil(t, model.ActorIP)
	assert.Equal(t, "127.0.0.1", *model.ActorIP)
}

func TestTableName(t *testing.T) {
	model := auditEventDB{}
	assert.Equal(t, "audit.events", model.TableName())
}
