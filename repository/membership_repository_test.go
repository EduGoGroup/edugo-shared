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
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	membership := &entities.Membership{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		SchoolID:       uuid.New(),
		Role:           "student",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectBegin()
	// GORM inserts all fields.
	// The table name is "academic"."memberships".
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "academic"."memberships"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), membership)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByID(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id", "role", "is_active"}).
		AddRow(id, uuid.New(), "student", true)

	// GORM adds ORDER BY "memberships"."id" LIMIT 1 by default for First()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	membership, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if membership != nil {
		assert.Equal(t, id, membership.ID)
	} else {
		t.Error("Membership should not be nil")
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUser(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	userID := uuid.New()
	filters := ListFilters{}

	rows := sqlmock.NewRows([]string{"id", "user_id", "role", "is_active"}).
		AddRow(uuid.New(), userID, "student", true).
		AddRow(uuid.New(), userID, "teacher", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE (user_id = $1 AND is_active = true) ORDER BY created_at DESC`)).
		WithArgs(userID).
		WillReturnRows(rows)

	memberships, err := repo.FindByUser(context.Background(), userID, filters)
	assert.NoError(t, err)
	assert.Len(t, memberships, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnit(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	unitID := uuid.New()
	filters := ListFilters{}

	rows := sqlmock.NewRows([]string{"id", "academic_unit_id", "role", "is_active"}).
		AddRow(uuid.New(), unitID, "student", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND is_active = true) ORDER BY created_at DESC`)).
		WithArgs(unitID).
		WillReturnRows(rows)

	memberships, err := repo.FindByUnit(context.Background(), unitID, filters)
	assert.NoError(t, err)
	assert.Len(t, memberships, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUnitAndRole(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	unitID := uuid.New()
	role := "teacher"
	activeOnly := true
	filters := ListFilters{}

	rows := sqlmock.NewRows([]string{"id", "academic_unit_id", "role", "is_active"}).
		AddRow(uuid.New(), unitID, role, true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND role = ?) AND is_active = true`)).
		WithArgs(unitID, role).
		WillReturnRows(rows)

	memberships, err := repo.FindByUnitAndRole(context.Background(), unitID, role, activeOnly, filters)
	assert.NoError(t, err)
	assert.Len(t, memberships, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_FindByUserAndSchool(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	userID := uuid.New()
	schoolID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id", "school_id", "is_active"}).
		AddRow(uuid.New(), userID, schoolID, true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE (user_id = $1 AND school_id = $2 AND is_active = true) ORDER BY "memberships"."id" LIMIT $3`)).
		WithArgs(userID, schoolID, 1).
		WillReturnRows(rows)

	membership, err := repo.FindByUserAndSchool(context.Background(), userID, schoolID)
	assert.NoError(t, err)
	assert.NotNil(t, membership)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Update(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	membership := &entities.Membership{
		ID:        uuid.New(),
		Role:      "admin",
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."memberships"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), membership)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMembershipRepository_Delete(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	// Hard delete expected for Membership as it probably has no DeletedAt
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "academic"."memberships" WHERE "memberships"."id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestMembershipRepository_NotFound ensures that ErrRecordNotFound is returned correctly
func TestMembershipRepository_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresMembershipRepository(db)
	id := uuid.New()

	// Adjusted for GORM First() default ordering and limit
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	membership, err := repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, membership)

	assert.NoError(t, mock.ExpectationsWereMet())
}
