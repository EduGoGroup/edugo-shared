package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSchoolRepository_Create(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	school := &entities.School{
		ID:        uuid.New(),
		Name:      "Test School",
		Code:      "TS001",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	// Matches INSERT INTO "academic"."schools" ...
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "academic"."schools"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), school)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByID(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "name", "code", "is_active"}).
		AddRow(id, "Test School", "TS001", true)

	// Handles "academic"."schools", "deleted_at" check, and LIMIT/ORDER BY
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	school, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if school != nil {
		assert.Equal(t, id, school.ID)
	} else {
		t.Error("School should not be nil")
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByCode(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	code := "TS001"
	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "name", "code", "is_active"}).
		AddRow(id, "Test School", code, true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(code, 1).
		WillReturnRows(rows)

	school, err := repo.FindByCode(context.Background(), code)
	assert.NoError(t, err)
	if school != nil {
		assert.Equal(t, code, school.Code)
	} else {
		t.Error("School should not be nil")
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_ExistsByCode(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	code := "TS001"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL`)).
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.ExistsByCode(context.Background(), code)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Update(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	school := &entities.School{
		ID:        uuid.New(),
		Name:      "Updated School",
		Code:      "TS001",
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	// Soft delete filtering in WHERE clause
	// Fixed: removed escaping backslashes that were causing issues in previous run
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), school)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Delete(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	// Soft delete expected
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET "deleted_at"=$1 WHERE id = $2 AND "schools"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_List(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	filters := ListFilters{
		Limit:  10,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "code"}).
		AddRow(uuid.New(), "School 1", "S001").
		AddRow(uuid.New(), "School 2", "S002")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE "schools"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1`)).
		WithArgs(10).
		WillReturnRows(rows)

	schools, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, schools, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSchoolRepository_NotFound ensures that ErrRecordNotFound is returned correctly
func TestSchoolRepository_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresSchoolRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	school, err := repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, school)

	assert.NoError(t, mock.ExpectationsWereMet())
}
