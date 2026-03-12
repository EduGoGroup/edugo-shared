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
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	return gormDB, mock
}

func TestNewPostgresUserRepository(t *testing.T) {
	gormDB, _ := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)
	assert.NotNil(t, repo)
}

func TestUserRepository_Create(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "test@test.com",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "auth"."users" ("id","email","password_hash","first_name","last_name","is_active","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).
		WithArgs(
			user.ID,
			user.Email,
			sqlmock.AnyArg(),
			user.FirstName,
			user.LastName,
			user.IsActive,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)
	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "is_active"}).
		AddRow(id, "test@test.com", "Test", "User", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	user, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	if user != nil {
		assert.Equal(t, id, user.ID)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)
	email := "test@test.com"

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "is_active"}).
		AddRow(uuid.New(), email, "Test", "User", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	if user != nil {
		assert.Equal(t, email, user.Email)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)
	email := "test@test.com"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(email).
		WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "test@test.com",
		FirstName: "Test",
		LastName:  "User Updated",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET "email"=$1,"password_hash"=$2,"first_name"=$3,"last_name"=$4,"is_active"=$5,"created_at"=$6,"updated_at"=$7,"deleted_at"=$8 WHERE "users"."deleted_at" IS NULL AND "id" = $9`)).
		WithArgs(
			user.Email,
			sqlmock.AnyArg(),
			user.FirstName,
			user.LastName,
			user.IsActive,
			user.CreatedAt,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			user.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)
	id := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_List(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewPostgresUserRepository(gormDB)

	isActive := true
	filters := ListFilters{
		IsActive: &isActive,
		Search:   "test",
		SearchFields: []string{"email"},
		Limit:    10,
		Offset:   0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE is_active = $1 AND email ILIKE $2 ESCAPE '\' AND "users"."deleted_at" IS NULL`)).
		WithArgs(isActive, "%test%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "is_active"}).
		AddRow(uuid.New(), "test1@test.com", "Test1", "User", true).
		AddRow(uuid.New(), "test2@test.com", "Test2", "User", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE is_active = $1 AND email ILIKE $2 ESCAPE '\' AND "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $3`)).
		WithArgs(isActive, "%test%", 10).
		WillReturnRows(rows)

	users, total, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
