package postgres

import "gorm.io/gorm"

// GORMTxFunc is a function that runs within a GORM transaction.
type GORMTxFunc func(*gorm.DB) error

// WithGORMTransaction executes the given function within a GORM transaction.
// If the function returns an error, the transaction is rolled back.
// If the function panics, the transaction is rolled back and the panic is re-raised.
func WithGORMTransaction(db *gorm.DB, fn GORMTxFunc) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
