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

func TestUserRepository_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	id := uuid.New()
	u := &entities.User{
		ID:       id,
		Email:    "test@example.com",
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "auth"\."users"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), u)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	id := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "is_active"}).
			AddRow(id.String(), "test@example.com", true))

	u, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if assert.NotNil(t, u) {
		assert.Equal(t, id, u.ID)
		assert.Equal(t, "test@example.com", u.Email)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	id := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	u, err := repo.FindByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, u)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	email := "test@example.com"

	mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(uuid.New().String(), email))

	u, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	if assert.NotNil(t, u) {
		assert.Equal(t, email, u.Email)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	id := uuid.New()
	u := &entities.User{
		ID:       id,
		Email:    "test2@example.com",
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "auth"\."users" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), u)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "auth"\."users" SET "deleted_at"=\$1 WHERE id = \$2 AND "users"\."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_List(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	isActive := true
	filters := ListFilters{
		Limit:    10,
		Offset:   0,
		IsActive: &isActive,
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE is_active = \$1 AND "users"."deleted_at" IS NULL`).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE is_active = \$1 AND "users"\."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$2`).
		WithArgs(true, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(uuid.New().String(), "test1@example.com").
			AddRow(uuid.New().String(), "test2@example.com"))

	users, total, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(db)

	email := "test@example.com"

	mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE email = \$1 AND "users"."deleted_at" IS NULL`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.ExistsByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
