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

func TestMembershipRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	unitID := uuid.New()
	schoolID := uuid.New()
	membership := &entities.Membership{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		SchoolID:       schoolID,
		AcademicUnitID: &unitID,
		Role:           "student",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."memberships"`)).
		WithArgs(membership.ID, membership.UserID, membership.SchoolID, membership.AcademicUnitID, membership.Role, membership.IsActive, membership.EnrolledAt, membership.CreatedAt, membership.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"metadata", "withdrawn_at"}).AddRow([]byte{}, nil))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), membership)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).AddRow(id, "student"))

	membership, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if assert.NotNil(t, membership) {
		assert.Equal(t, id, membership.ID)
		assert.Equal(t, "student", membership.Role)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{}))

	membership, err := repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, membership)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	userID := uuid.New()

	filters := ListFilters{
		Limit:  10,
		Offset: 0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE user_id = $1 AND is_active = true`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC LIMIT $2`)).
		WithArgs(userID, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.New(), "student").
			AddRow(uuid.New(), "teacher"))

	memberships, total, err := repo.FindByUser(context.Background(), userID, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, memberships, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnit(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	unitID := uuid.New()

	filters := ListFilters{
		Limit:  10,
		Offset: 0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE academic_unit_id = $1 AND is_active = true`)).
		WithArgs(unitID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE academic_unit_id = $1 AND is_active = true ORDER BY created_at DESC LIMIT $2`)).
		WithArgs(unitID, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.New(), "student").
			AddRow(uuid.New(), "teacher"))

	memberships, total, err := repo.FindByUnit(context.Background(), unitID, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, memberships, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnitAndRole(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	unitID := uuid.New()
	role := "student"

	filters := ListFilters{
		Limit:  10,
		Offset: 0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND role = $2) AND is_active = true`)).
		WithArgs(unitID, role).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND role = $2) AND is_active = true ORDER BY created_at DESC LIMIT $3`)).
		WithArgs(unitID, role, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.New(), role).
			AddRow(uuid.New(), role))

	memberships, total, err := repo.FindByUnitAndRole(context.Background(), unitID, role, true, filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, memberships, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUserAndSchool(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	userID := uuid.New()
	schoolID := uuid.New()
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND school_id = $2 AND is_active = true ORDER BY "memberships"."id" LIMIT $3`)).
		WithArgs(userID, schoolID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).AddRow(id, "student"))

	membership, err := repo.FindByUserAndSchool(context.Background(), userID, schoolID)
	assert.NoError(t, err)
	if assert.NotNil(t, membership) {
		assert.Equal(t, id, membership.ID)
		assert.Equal(t, "student", membership.Role)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	membership := &entities.Membership{
		ID:        uuid.New(),
		Role:      "teacher",
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."memberships" SET`)).
		WithArgs(membership.UserID, membership.SchoolID, membership.AcademicUnitID, membership.Role, membership.Metadata, membership.IsActive, membership.EnrolledAt, membership.WithdrawnAt, sqlmock.AnyArg(), sqlmock.AnyArg(), membership.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), membership)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "academic"."memberships" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
