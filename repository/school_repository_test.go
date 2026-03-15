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

func setupSchoolMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestSchoolRepository_Create(t *testing.T) {
	gormDB, mock := setupSchoolMockDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	school := &entities.School{
		ID:        uuid.New(),
		Code:      "SCH001",
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
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."schools"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(school.ID))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."schools"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Create(context.Background(), school)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSchoolRepository_FindByID(t *testing.T) {
	gormDB, mock := setupSchoolMockDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	schoolID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    *entities.School
		wantErr error
	}{
		{
			name: "Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "code"}).
					AddRow(schoolID, "SCH001")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
					WithArgs(schoolID, 1).
					WillReturnRows(rows)
			},
			want: &entities.School{
				ID:   schoolID,
				Code: "SCH001",
			},
			wantErr: nil,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
					WithArgs(schoolID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: ErrNotFound,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
					WithArgs(schoolID, 1).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindByID(context.Background(), schoolID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, ErrNotFound) {
					assert.Equal(t, ErrNotFound, err)
				}
			} else {
				assert.NoError(t, err)
				if got != nil && tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.Code, got.Code)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSchoolRepository_FindByCode(t *testing.T) {
	gormDB, mock := setupSchoolMockDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	code := "SCH001"
	schoolID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    *entities.School
		wantErr error
	}{
		{
			name: "Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "code"}).
					AddRow(schoolID, code)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
					WithArgs(code, 1).
					WillReturnRows(rows)
			},
			want: &entities.School{
				ID:   schoolID,
				Code: code,
			},
			wantErr: nil,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
					WithArgs(code, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: ErrNotFound,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
					WithArgs(code, 1).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindByCode(context.Background(), code)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, ErrNotFound) {
					assert.Equal(t, ErrNotFound, err)
				}
			} else {
				assert.NoError(t, err)
				if got != nil && tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.Code, got.Code)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSchoolRepository_Update(t *testing.T) {
	gormDB, mock := setupSchoolMockDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	school := &entities.School{
		ID:   uuid.New(),
		Code: "SCH001",
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
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Update(context.Background(), school)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSchoolRepository_Delete(t *testing.T) {
	gormDB, mock := setupSchoolMockDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	schoolID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET "deleted_at"=$1 WHERE id = $2 AND "schools"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), schoolID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET "deleted_at"=$1 WHERE id = $2 AND "schools"."deleted_at" IS NULL`)).
					WithArgs(sqlmock.AnyArg(), schoolID).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Delete(context.Background(), schoolID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSchoolRepository_List(t *testing.T) {
	gormDB, mock := setupSchoolMockDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	schoolID := uuid.New()

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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL`)).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "code"}).AddRow(schoolID, "SCH001")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL ORDER BY created_at DESC`)).
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE is_active = $1 AND "schools"."deleted_at" IS NULL`)).
					WithArgs(true).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "code"}).AddRow(schoolID, "SCH001")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE is_active = $1 AND "schools"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(true).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "With Search",
			filters: ListFilters{
				Search:       "SCH001",
				SearchFields: []string{"code"},
			},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE code ILIKE $1 ESCAPE '\' AND "schools"."deleted_at" IS NULL`)).
					WithArgs("%SCH001%").
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "code"}).AddRow(schoolID, "SCH001")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE code ILIKE $1 ESCAPE '\' AND "schools"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs("%SCH001%").
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL`)).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "code"}).AddRow(schoolID, "SCH001")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL`)).
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL`)).
					WillReturnRows(countRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WillReturnError(errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			schools, total, err := repo.List(context.Background(), tt.filters)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, schools, tt.wantLen)
				if tt.wantLen > 0 {
					assert.Equal(t, int64(1), total)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSchoolRepository_ExistsByCode(t *testing.T) {
	gormDB, mock := setupSchoolMockDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	code := "SCH001"

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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL`)).
					WithArgs(code).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Does Not Exist",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL`)).
					WithArgs(code).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL`)).
					WithArgs(code).
					WillReturnError(errors.New("db error"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.ExistsByCode(context.Background(), code)
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
