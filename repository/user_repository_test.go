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

func TestUserRepository_Create(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresUserRepository(db)

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	// Matches INSERT INTO "auth"."users" ...
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "auth"."users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresUserRepository(db)

	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "is_active"}).
		AddRow(id, "test@example.com", "John", "Doe", true)

	// Matches SELECT ... FROM "auth"."users" ... with deleted_at check and limit
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	user, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if user != nil {
		assert.Equal(t, id, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
	} else {
		t.Error("User should not be nil")
	}

	// Not found case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err = repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresUserRepository(db)

	email := "test@example.com"
	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "is_active"}).
		AddRow(id, email, "John", "Doe", true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	if user != nil {
		assert.Equal(t, email, user.Email)
	} else {
		t.Error("User should not be nil")
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresUserRepository(db)

	email := "test@example.com"

	// Exists
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.ExistsByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Does not exist
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "auth"."users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err = repo.ExistsByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPostgresUserRepository(db)

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "updated@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth"."users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock := newMockDB(t)
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
	db, mock := newMockDB(t)
	repo := NewPostgresUserRepository(db)

	// Case 1: Simple list with limit and offset
	filters := ListFilters{
		Limit:  10,
		Offset: 0,
	}

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(uuid.New(), "user1@example.com").
		AddRow(uuid.New(), "user2@example.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1`)).
		WithArgs(10).
		WillReturnRows(rows)

	users, err := repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Case 2: Filter by IsActive
	isActive := true
	filters = ListFilters{
		IsActive: &isActive,
	}

	// Reuse rows
	rows2 := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(uuid.New(), "user1@example.com").
		AddRow(uuid.New(), "user2@example.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE is_active = $1 AND "users"."deleted_at" IS NULL ORDER BY created_at DESC`)).
		WithArgs(true).
		WillReturnRows(rows2)

	users, err = repo.List(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Case 3: Search
	filters = ListFilters{
		Search:       "John",
		SearchFields: []string{"first_name", "last_name"},
	}

	rows3 := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(uuid.New(), "user1@example.com").
		AddRow(uuid.New(), "user2@example.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth"."users" WHERE (first_name ILIKE $1 ESCAPE '\' OR last_name ILIKE $2 ESCAPE '\') AND "users"."deleted_at" IS NULL ORDER BY created_at DESC`)).
		WithArgs("%John%", "%John%").
		WillReturnRows(rows3)

	users, err = repo.List(context.Background(), filters)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
