package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

// TxFunc es una función que se ejecuta dentro de una transacción
type TxFunc func(*sql.Tx) error

// WithTransaction ejecuta una función dentro de una transacción
// Si la función retorna error, hace rollback. Si no, hace commit.
func WithTransaction(ctx context.Context, db *sql.DB, fn TxFunc) error {
	// Comenzar transacción
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer para manejar rollback en caso de panic
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic después de rollback
		}
	}()

	// Ejecutar función
	if err := fn(tx); err != nil {
		// Rollback en caso de error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	// Commit si todo salió bien
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionIsolation ejecuta una función dentro de una transacción con nivel de aislamiento específico
func WithTransactionIsolation(ctx context.Context, db *sql.DB, isolation sql.IsolationLevel, fn TxFunc) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: isolation,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
