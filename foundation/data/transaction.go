package data

import (
	"gorm.io/gorm"
)

// TxFn  A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(db *gorm.DB) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(db *gorm.DB, fn TxFn) (err error) {
	tx := db.Begin()

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
