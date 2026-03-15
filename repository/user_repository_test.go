package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupUserMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestUserRepository_Create(t *testing.T) {
	gormDB, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(gormDB)

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "auth"."users"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "auth"."users"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Create(context.Background(), user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	gormDB, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(gormDB)

	userID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    *entities.User
		wantErr error
	}{
		{
			name: "Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(userID, "test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnRows(rows)
			},
			want: &entities.User{
				ID:    userID,
				Email: "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: ErrNotFound,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindByID(context.Background(), userID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, ErrNotFound) {
					assert.Equal(t, ErrNotFound, err)
				}
			} else {
				assert.NoError(t, err)
				if got != nil && tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.Email, got.Email)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	gormDB, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(gormDB)

	email := "test@example.com"
	userID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    *entities.User
		wantErr error
	}{
		{
			name: "Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(userID, email)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(email, 1).
					WillReturnRows(rows)
			},
			want: &entities.User{
				ID:    userID,
				Email: email,
			},
			wantErr: nil,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(email, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: ErrNotFound,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(email, 1).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindByEmail(context.Background(), email)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, ErrNotFound) {
					assert.Equal(t, ErrNotFound, err)
				}
			} else {
				assert.NoError(t, err)
				if got != nil && tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.Email, got.Email)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	gormDB, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(gormDB)

	email := "test@example.com"

	tests := []struct {
		name    string
		mock    func()
		want    bool
		wantErr bool
	}{
		{
			name: "Exists",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
					WithArgs(email).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Does Not Exist",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
					WithArgs(email).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
					WithArgs(email).
					WillReturnError(errors.New("db error"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.ExistsByEmail(context.Background(), email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	gormDB, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(gormDB)

	user := &entities.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Update(context.Background(), user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	gormDB, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(gormDB)

	userID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), userID).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Delete(context.Background(), userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_List(t *testing.T) {
	gormDB, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(gormDB)

	userID := uuid.New()

	tests := []struct {
		name    string
		filters ListFilters
		mock    func()
		wantLen int
		wantErr bool
	}{
		{
			name:    "No Filters",
			filters: ListFilters{},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE "users"."deleted_at" IS NULL`)).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, "test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE "users"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "With Active Filter",
			filters: ListFilters{
				IsActive: func(b bool) *bool { return &b }(true),
			},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE is_active = $1 AND "users"."deleted_at" IS NULL`)).
					WithArgs(true).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, "test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE is_active = $1 AND "users"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(true).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "With Search",
			filters: ListFilters{
				Search:       "test",
				SearchFields: []string{"email"},
			},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email ILIKE $1 ESCAPE '\' AND "users"."deleted_at" IS NULL`)).
					WithArgs("%test%").
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, "test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE email ILIKE $1 ESCAPE '\' AND "users"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs("%test%").
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "With Pagination",
			filters: ListFilters{
				Limit:  10,
				Offset: 5,
			},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE "users"."deleted_at" IS NULL`)).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, "test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
					WithArgs(10, 5).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "Count Error",
			filters: ListFilters{},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE "users"."deleted_at" IS NULL`)).
					WillReturnError(errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "Find Error",
			filters: ListFilters{},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE "users"."deleted_at" IS NULL`)).
					WillReturnRows(countRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE "users"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WillReturnError(errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			users, total, err := repo.List(context.Background(), tt.filters)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, tt.wantLen)
				if tt.wantLen > 0 {
					assert.Equal(t, int64(1), total)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
