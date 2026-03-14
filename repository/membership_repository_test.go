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

func TestMembershipRepository_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	id := uuid.New()
	userID := uuid.New()
	schoolID := uuid.New()
	unitID := uuid.New()

	m := &entities.Membership{
		ID:             id,
		UserID:         userID,
		SchoolID:       schoolID,
		AcademicUnitID: &unitID,
		Role:           "student",
		IsActive:       true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "academic"."memberships"`).
		WillReturnRows(sqlmock.NewRows([]string{"metadata", "withdrawn_at"}).AddRow(nil, nil))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), m)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByID(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	id := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "academic"\."memberships" WHERE id = \$1 ORDER BY "memberships"."id" LIMIT \$2`).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role", "is_active"}).
			AddRow(id.String(), "student", true))

	m, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if assert.NotNil(t, m) {
		assert.Equal(t, id, m.ID)
		assert.Equal(t, "student", m.Role)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	id := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "academic"\."memberships" WHERE id = \$1 ORDER BY "memberships"."id" LIMIT \$2`).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	m, err := repo.FindByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, m)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUser(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	userID := uuid.New()
	filters := ListFilters{Limit: 10, Offset: 0}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "academic"\."memberships" WHERE user_id = \$1 AND is_active = true`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(`SELECT \* FROM "academic"\."memberships" WHERE user_id = \$1 AND is_active = true ORDER BY created_at DESC LIMIT \$2`).
		WithArgs(userID, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).
			AddRow(uuid.New().String(), userID.String()).
			AddRow(uuid.New().String(), userID.String()))

	memberships, total, err := repo.FindByUser(context.Background(), userID, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, memberships, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnit(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	unitID := uuid.New()
	filters := ListFilters{Limit: 10, Offset: 0}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "academic"\."memberships" WHERE academic_unit_id = \$1 AND is_active = true`).
		WithArgs(unitID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT \* FROM "academic"\."memberships" WHERE academic_unit_id = \$1 AND is_active = true ORDER BY created_at DESC LIMIT \$2`).
		WithArgs(unitID, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "academic_unit_id"}).
			AddRow(uuid.New().String(), unitID.String()))

	memberships, total, err := repo.FindByUnit(context.Background(), unitID, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, memberships, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnitAndRole(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	unitID := uuid.New()
	filters := ListFilters{Limit: 10, Offset: 0}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "academic"\."memberships" WHERE \(academic_unit_id = \$1 AND role = \$2\) AND is_active = true`).
		WithArgs(unitID, "student").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT \* FROM "academic"\."memberships" WHERE \(academic_unit_id = \$1 AND role = \$2\) AND is_active = true ORDER BY created_at DESC LIMIT \$3`).
		WithArgs(unitID, "student", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "academic_unit_id", "role"}).
			AddRow(uuid.New().String(), unitID.String(), "student"))

	memberships, total, err := repo.FindByUnitAndRole(context.Background(), unitID, "student", true, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, memberships, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUserAndSchool(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	userID := uuid.New()
	schoolID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "academic"\."memberships" WHERE user_id = \$1 AND school_id = \$2 AND is_active = true ORDER BY "memberships"."id" LIMIT \$3`).
		WithArgs(userID, schoolID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "school_id"}).
			AddRow(uuid.New().String(), userID.String(), schoolID.String()))

	m, err := repo.FindByUserAndSchool(context.Background(), userID, schoolID)
	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	id := uuid.New()
	m := &entities.Membership{
		ID:       id,
		Role:     "teacher",
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "academic"."memberships" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), m)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresMembershipRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "academic"\."memberships" WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
