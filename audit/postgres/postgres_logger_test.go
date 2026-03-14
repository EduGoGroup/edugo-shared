package postgres

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
	}

	return db, mock
}

func TestPostgresAuditLogger_Log(t *testing.T) {
	tests := []struct {
		name          string
		event         audit.AuditEvent
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "successful log with defaults",
			event: audit.AuditEvent{
				ActorID:      "user-123",
				ActorEmail:   "test@example.com",
				ActorRole:    "admin",
				Action:       "CREATE",
				ResourceType: "USER",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "audit"."events"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("test-uuid"))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "database error",
			event: audit.AuditEvent{
				ActorID:      "user-123",
				ActorEmail:   "test@example.com",
				ActorRole:    "admin",
				Action:       "CREATE",
				ResourceType: "USER",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "audit"."events"`).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			logger := NewPostgresAuditLogger(db, "test-service")

			tt.mockSetup(mock)

			err := logger.Log(context.Background(), tt.event)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresAuditLogger_LogFromGin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		setupContext  func(c *gin.Context)
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "successful log from gin context",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "user-123")
				c.Set("email", "test@example.com")
				c.Set("role", "admin")
				c.Request.Header.Set("X-Request-ID", "req-123")
				c.Request.Header.Set("User-Agent", "test-agent")
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "audit"."events"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("test-uuid"))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupTestDB(t)
			logger := NewPostgresAuditLogger(db, "test-service")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
			c.Request.RemoteAddr = "127.0.0.1:12345"

			tt.setupContext(c)
			tt.mockSetup(mock)

			err := logger.LogFromGin(c, "CREATE", "USER", "res-123")

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
