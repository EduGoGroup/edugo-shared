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
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	school := &entities.School{
		ID:        uuid.New(),
		Name:      "Test School",
		Code:      "TS-001",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."schools"`)).
		WithArgs(school.ID, school.Name, school.Code, school.Country, school.ConceptTypeID, school.IsActive, school.SubscriptionTier, school.MaxTeachers, school.MaxStudents, school.CreatedAt, school.UpdatedAt, school.DeletedAt).
		WillReturnRows(sqlmock.NewRows([]string{"address", "city", "phone", "email", "metadata"}).AddRow("", "", "", "", []byte{}))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), school)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).AddRow(id, "TS-001", "Test School"))

	school, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if assert.NotNil(t, school) {
		assert.Equal(t, id, school.ID)
		assert.Equal(t, "TS-001", school.Code)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{}))

	school, err := repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, school)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByCode(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)
	code := "TS-001"
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(code, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).AddRow(id, code, "Test School"))

	school, err := repo.FindByCode(context.Background(), code)
	assert.NoError(t, err)
	if assert.NotNil(t, school) {
		assert.Equal(t, code, school.Code)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)
	school := &entities.School{
		ID:        uuid.New(),
		Name:      "Updated School",
		Code:      "TS-002",
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET`)).
		WithArgs(school.Name, school.Code, school.Address, school.City, school.Country, school.Phone, school.Email, school.ConceptTypeID, school.Metadata, school.IsActive, school.SubscriptionTier, school.MaxTeachers, school.MaxStudents, sqlmock.AnyArg(), sqlmock.AnyArg(), school.DeletedAt, school.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), school)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET "deleted_at"=$1 WHERE id = $2 AND "schools"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_List(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)

	isActive := true
	filters := ListFilters{
		IsActive: &isActive,
		Limit:    10,
		Offset:   0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE is_active = $1 AND "schools"."deleted_at" IS NULL`)).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE is_active = $1 AND "schools"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $2`)).
		WithArgs(true, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).
			AddRow(uuid.New(), "TS-001", "School 1").
			AddRow(uuid.New(), "TS-002", "School 2"))

	schools, total, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, schools, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_ExistsByCode(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresSchoolRepository(db)
	code := "TS-001"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL`)).
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.ExistsByCode(context.Background(), code)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
