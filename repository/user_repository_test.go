package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupUserMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	return db, mock
}

func TestUserRepository_Create(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	id := uuid.New()
	user := &entities.User{
		ID:        id,
		Email:     "test@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "auth"\."users"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(id, "test@example.com"))

		user, err := repo.FindByID(ctx, id)
		assert.NoError(t, err)
		if assert.NotNil(t, user) {
			assert.Equal(t, id, user.ID)
			assert.Equal(t, "test@example.com", user.Email)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

		user, err := repo.FindByID(ctx, id)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		dbErr := errors.New("db error")
		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE id = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(id, 1).
			WillReturnError(dbErr)

		user, err := repo.FindByID(ctx, id)
		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()
	email := "test@example.com"

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(uuid.New(), email))

		user, err := repo.FindByEmail(ctx, email)
		assert.NoError(t, err)
		if assert.NotNil(t, user) {
			assert.Equal(t, email, user.Email)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

		user, err := repo.FindByEmail(ctx, email)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		dbErr := errors.New("db error")
		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnError(dbErr)

		user, err := repo.FindByEmail(ctx, email)
		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()
	email := "test@example.com"

	t.Run("exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		exists, err := repo.ExistsByEmail(ctx, email)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		exists, err := repo.ExistsByEmail(ctx, email)
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		dbErr := errors.New("db error")
		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL`).
			WithArgs(email).
			WillReturnError(dbErr)

		exists, err := repo.ExistsByEmail(ctx, email)
		assert.ErrorIs(t, err, dbErr)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Update(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	id := uuid.New()
	user := &entities.User{
		ID:        id,
		Email:     "updated@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "auth"\."users"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "auth"\."users" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_List(t *testing.T) {
	db, mock := setupUserMockDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	t.Run("success_without_filters", func(t *testing.T) {
		filters := ListFilters{}
		id1 := uuid.New()
		id2 := uuid.New()

		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE "users"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE "users"\."deleted_at" IS NULL ORDER BY created_at DESC`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
				AddRow(id1, "user1@example.com").
				AddRow(id2, "user2@example.com"))

		users, total, err := repo.List(ctx, filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, users, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_with_active_filter", func(t *testing.T) {
		isActive := true
		filters := ListFilters{IsActive: &isActive}
		id1 := uuid.New()

		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE is_active = \$1 AND "users"\."deleted_at" IS NULL`).
			WithArgs(isActive).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE is_active = \$1 AND "users"\."deleted_at" IS NULL ORDER BY created_at DESC`).
			WithArgs(isActive).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(id1, "active@example.com"))

		users, total, err := repo.List(ctx, filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, users, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_with_pagination", func(t *testing.T) {
		filters := ListFilters{Limit: 10, Offset: 20}
		id1 := uuid.New()

		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE "users"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE "users"\."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
			WithArgs(10, 20).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(id1, "user@example.com"))

		users, total, err := repo.List(ctx, filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(50), total)
		assert.Len(t, users, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("count_error", func(t *testing.T) {
		filters := ListFilters{}
		dbErr := errors.New("db error")

		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE "users"\."deleted_at" IS NULL`).
			WillReturnError(dbErr)

		users, total, err := repo.List(ctx, filters)
		assert.ErrorIs(t, err, dbErr)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("find_error", func(t *testing.T) {
		filters := ListFilters{}
		dbErr := errors.New("db error")

		mock.ExpectQuery(`SELECT count\(\*\) FROM "auth"\."users" WHERE "users"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(`SELECT \* FROM "auth"\."users" WHERE "users"\."deleted_at" IS NULL ORDER BY created_at DESC`).
			WillReturnError(dbErr)

		users, total, err := repo.List(ctx, filters)
		assert.ErrorIs(t, err, dbErr)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
