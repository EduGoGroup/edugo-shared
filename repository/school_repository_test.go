package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSchoolRepository_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	id := uuid.New()
	s := &entities.School{
		ID:       id,
		Code:     "SCHOOL-01",
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "academic"\."schools"`).
		WillReturnRows(sqlmock.NewRows([]string{"metadata"}).AddRow(nil))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), s)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByID(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	id := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "academic"\."schools" WHERE id = \$1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT \$2`).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "is_active"}).
			AddRow(id.String(), "SCHOOL-01", true))

	s, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if assert.NotNil(t, s) {
		assert.Equal(t, id, s.ID)
		assert.Equal(t, "SCHOOL-01", s.Code)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	id := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "academic"\."schools" WHERE id = \$1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT \$2`).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	s, err := repo.FindByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByCode(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	code := "SCHOOL-01"

	mock.ExpectQuery(`SELECT \* FROM "academic"\."schools" WHERE code = \$1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT \$2`).
		WithArgs(code, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code"}).
			AddRow(uuid.New().String(), code))

	s, err := repo.FindByCode(context.Background(), code)
	assert.NoError(t, err)
	if assert.NotNil(t, s) {
		assert.Equal(t, code, s.Code)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	id := uuid.New()
	s := &entities.School{
		ID:       id,
		Code:     "SCHOOL-02",
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "academic"\."schools" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), s)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "academic"\."schools" SET "deleted_at"=\$1 WHERE id = \$2 AND "schools"\."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_List(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	isActive := true
	filters := ListFilters{
		Limit:    10,
		Offset:   0,
		IsActive: &isActive,
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "academic"\."schools" WHERE is_active = \$1 AND "schools"."deleted_at" IS NULL`).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(`SELECT \* FROM "academic"\."schools" WHERE is_active = \$1 AND "schools"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$2`).
		WithArgs(true, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code"}).
			AddRow(uuid.New().String(), "S1").
			AddRow(uuid.New().String(), "S2"))

	schools, total, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, schools, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_ExistsByCode(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(db)

	code := "SCHOOL-01"

	mock.ExpectQuery(`SELECT count\(\*\) FROM "academic"\."schools" WHERE code = \$1 AND "schools"."deleted_at" IS NULL`).
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.ExistsByCode(context.Background(), code)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
