package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// TxFunc es una funcion que se ejecuta dentro de una transaccion SQL.
type TxFunc func(*sql.Tx) error

// GORMTxFunc es una funcion que se ejecuta dentro de una transaccion GORM.
type GORMTxFunc func(*gorm.DB) error

// WithTransaction ejecuta una funcion dentro de una transaccion.
// Si la funcion retorna error, hace rollback. Si no, hace commit.
// Si ocurre un panic, hace rollback y re-lanza el panic.
func WithTransaction(ctx context.Context, db *sql.DB, fn TxFunc) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() //nolint:errcheck // En panic no hay forma practica de manejar error de rollback
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(
				fmt.Errorf("tx error: %w", err),
				fmt.Errorf("rollback error: %w", rbErr),
			)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionIsolation ejecuta una funcion dentro de una transaccion
// con nivel de aislamiento especifico.
func WithTransactionIsolation(ctx context.Context, db *sql.DB, isolation sql.IsolationLevel, fn TxFunc) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: isolation,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() //nolint:errcheck // En panic no hay forma practica de manejar error de rollback
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(
				fmt.Errorf("tx error: %w", err),
				fmt.Errorf("rollback error: %w", rbErr),
			)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithGORMTransaction ejecuta una funcion dentro de una transaccion GORM.
// Si la funcion retorna error, la transaccion hace rollback.
// Si ocurre un panic, la transaccion hace rollback y re-lanza el panic.
func WithGORMTransaction(db *gorm.DB, fn GORMTxFunc) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
