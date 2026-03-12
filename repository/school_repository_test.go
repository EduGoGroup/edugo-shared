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
)

func TestNewPostgresSchoolRepository(t *testing.T) {
	gormDB, _ := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)
	assert.NotNil(t, repo)
}

func TestSchoolRepository_Create(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	school := &entities.School{
		ID:        uuid.New(),
		Code:      "SCH001",
		Name:      "Test School",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."schools" ("id","name","code","country","concept_type_id","is_active","subscription_tier","max_teachers","max_students","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "address","city","phone","email","metadata"`)).
		WithArgs(
			school.ID,
			school.Name,
			school.Code,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			school.IsActive,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"address","city","phone","email","metadata"}).AddRow("","","","",""))

	err := repo.Create(context.Background(), school)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByID(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)
	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "code", "name", "is_active"}).
		AddRow(id, "SCH001", "Test School", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE id = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	school, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, school)
	if school != nil {
		assert.Equal(t, id, school.ID)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_FindByCode(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)
	code := "SCH001"

	rows := sqlmock.NewRows([]string{"id", "code", "name", "is_active"}).
		AddRow(uuid.New(), code, "Test School", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL ORDER BY "schools"."id" LIMIT $2`)).
		WithArgs(code, 1).
		WillReturnRows(rows)

	school, err := repo.FindByCode(context.Background(), code)
	assert.NoError(t, err)
	assert.NotNil(t, school)
	if school != nil {
		assert.Equal(t, code, school.Code)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_ExistsByCode(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)
	code := "SCH001"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE code = $1 AND "schools"."deleted_at" IS NULL`)).
		WithArgs(code).
		WillReturnRows(rows)

	exists, err := repo.ExistsByCode(context.Background(), code)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Update(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	school := &entities.School{
		ID:        uuid.New(),
		Code:      "SCH001",
		Name:      "Test School Updated",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET "name"=$1,"code"=$2,"address"=$3,"city"=$4,"country"=$5,"phone"=$6,"email"=$7,"concept_type_id"=$8,"metadata"=$9,"is_active"=$10,"subscription_tier"=$11,"max_teachers"=$12,"max_students"=$13,"created_at"=$14,"updated_at"=$15,"deleted_at"=$16 WHERE "schools"."deleted_at" IS NULL AND "id" = $17`)).
		WithArgs(
			school.Name,
			school.Code,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			school.IsActive,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			school.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), school)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_Delete(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)
	id := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."schools" SET "deleted_at"=$1 WHERE id = $2 AND "schools"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolRepository_List(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresSchoolRepository(gormDB)

	isActive := true
	filters := ListFilters{
		IsActive: &isActive,
		Search:   "test",
		SearchFields: []string{"code"},
		Limit:    10,
		Offset:   0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."schools" WHERE is_active = $1 AND code ILIKE $2 ESCAPE '\' AND "schools"."deleted_at" IS NULL`)).
		WithArgs(isActive, "%test%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	rows := sqlmock.NewRows([]string{"id", "code", "name", "is_active"}).
		AddRow(uuid.New(), "SCH001", "Test1", true).
		AddRow(uuid.New(), "SCH002", "Test2", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."schools" WHERE is_active = $1 AND code ILIKE $2 ESCAPE '\' AND "schools"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $3`)).
		WithArgs(isActive, "%test%", 10).
		WillReturnRows(rows)

	schools, total, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, schools, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
