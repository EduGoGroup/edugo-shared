package postgres

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupGormMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestWithGORMTransaction_Success(t *testing.T) {
	db, mock := setupGormMockDB(t)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := WithGORMTransaction(db, func(tx *gorm.DB) error {
		// La transacción se ejecuta correctamente
		return nil
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithGORMTransaction_ErrorRollback(t *testing.T) {
	db, mock := setupGormMockDB(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	expectedErr := errors.New("error provocado para hacer rollback")
	err := WithGORMTransaction(db, func(tx *gorm.DB) error {
		return expectedErr
	})

	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithGORMTransaction_PanicRollback(t *testing.T) {
	db, mock := setupGormMockDB(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	defer func() {
		r := recover()
		assert.NotNil(t, r, "La función debió hacer panic")
		assert.Equal(t, "panic intencional", r)
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	err := WithGORMTransaction(db, func(tx *gorm.DB) error {
		panic("panic intencional")
	})
	assert.Error(t, err) // This line won't execute due to panic, but satisfies errcheck
}
