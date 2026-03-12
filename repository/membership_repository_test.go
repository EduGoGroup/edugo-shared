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

func TestNewPostgresMembershipRepository(t *testing.T) {
	gormDB, _ := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)
	assert.NotNil(t, repo)
}

func TestMembershipRepository_Create(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	unitID := uuid.New()
	m := &entities.Membership{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		SchoolID:       uuid.New(),
		AcademicUnitID: &unitID,
		Role:           "student",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."memberships" ("id","user_id","school_id","academic_unit_id","role","is_active","enrolled_at","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "metadata","withdrawn_at"`)).
		WithArgs(
			m.ID,
			m.UserID,
			m.SchoolID,
			m.AcademicUnitID,
			m.Role,
			m.IsActive,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"metadata", "withdrawn_at"}).AddRow("", nil))

	err := repo.Create(context.Background(), m)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByID(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)
	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "role", "is_active"}).
		AddRow(id, "student", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	m, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, m)
	if m != nil {
		assert.Equal(t, id, m.ID)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUser(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)
	userID := uuid.New()

	filters := ListFilters{
		Limit:  10,
		Offset: 0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE user_id = $1 AND is_active = true`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	rows := sqlmock.NewRows([]string{"id", "user_id", "role"}).
		AddRow(uuid.New(), userID, "student")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC LIMIT $2`)).
		WithArgs(userID, 10).
		WillReturnRows(rows)

	memberships, total, err := repo.FindByUser(context.Background(), userID, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, memberships, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnit(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)
	unitID := uuid.New()

	filters := ListFilters{
		Limit:  10,
		Offset: 0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE academic_unit_id = $1 AND is_active = true`)).
		WithArgs(unitID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	rows := sqlmock.NewRows([]string{"id", "academic_unit_id", "role"}).
		AddRow(uuid.New(), unitID, "student")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE academic_unit_id = $1 AND is_active = true ORDER BY created_at DESC LIMIT $2`)).
		WithArgs(unitID, 10).
		WillReturnRows(rows)

	memberships, total, err := repo.FindByUnit(context.Background(), unitID, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, memberships, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnitAndRole(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)
	unitID := uuid.New()
	role := "student"

	filters := ListFilters{
		Limit:  10,
		Offset: 0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND role = $2) AND is_active = true`)).
		WithArgs(unitID, role).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	rows := sqlmock.NewRows([]string{"id", "academic_unit_id", "role"}).
		AddRow(uuid.New(), unitID, role)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND role = $2) AND is_active = true ORDER BY created_at DESC LIMIT $3`)).
		WithArgs(unitID, role, 10).
		WillReturnRows(rows)

	memberships, total, err := repo.FindByUnitAndRole(context.Background(), unitID, role, true, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, memberships, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUserAndSchool(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)
	userID := uuid.New()
	schoolID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id", "school_id", "role"}).
		AddRow(uuid.New(), userID, schoolID, "student")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND school_id = $2 AND is_active = true ORDER BY "memberships"."id" LIMIT $3`)).
		WithArgs(userID, schoolID, 1).
		WillReturnRows(rows)

	m, err := repo.FindByUserAndSchool(context.Background(), userID, schoolID)
	assert.NoError(t, err)
	assert.NotNil(t, m)
	if m != nil {
		assert.Equal(t, userID, m.UserID)
		assert.Equal(t, schoolID, m.SchoolID)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Update(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	unitID := uuid.New()
	m := &entities.Membership{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		SchoolID:       uuid.New(),
		AcademicUnitID: &unitID,
		Role:           "teacher",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."memberships" SET "user_id"=$1,"school_id"=$2,"academic_unit_id"=$3,"role"=$4,"metadata"=$5,"is_active"=$6,"enrolled_at"=$7,"withdrawn_at"=$8,"created_at"=$9,"updated_at"=$10 WHERE "id" = $11`)).
		WithArgs(
			m.UserID,
			m.SchoolID,
			m.AcademicUnitID,
			m.Role,
			sqlmock.AnyArg(),
			m.IsActive,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			m.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), m)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Delete(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(gormDB)
	id := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "academic"."memberships" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
