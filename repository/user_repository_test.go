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
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestUserRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "auth"."users"`)).
		WithArgs(user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt, user.DeletedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(id, "test@example.com"))

	user, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, id, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{}))

	user, err := repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)
	email := "test@example.com"
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(id, email))

	user, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, email, user.Email)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)
	email := "test@example.com"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.ExistsByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)
	user := &entities.User{
		ID:        uuid.New(),
		Email:     "update@example.com",
		FirstName: "Updated",
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET`)).
		WithArgs(user.Email, user.PasswordHash, user.FirstName, user.LastName, user.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg(), user.DeletedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_List(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostgresUserRepository(db)

	isActive := true
	filters := ListFilters{
		IsActive: &isActive,
		Limit:    10,
		Offset:   0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE is_active = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE is_active = $1 AND "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $2`)).
		WithArgs(true, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(uuid.New(), "user1@example.com").
			AddRow(uuid.New(), "user2@example.com"))

	users, total, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
