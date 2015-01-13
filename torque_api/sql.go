package torque

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// sql.go provides useful SQL functions.
// Mostly SQL templates.

// DBActor defines an object which implements basic data operations
type DBActor interface {
	Create(*sqlx.DB) error
	Retrieve(*sqlx.DB) error
	Update(*sqlx.DB) error
	Delete(*sqlx.DB) error
}

// Transactor adds operations to sql.Transactions
type Transactor interface {
	Create(*sqlx.Tx) error
	Retrieve(*sqlx.Tx) error
	Update(*sqlx.Tx) error
	Delete(*sqlx.Tx) error
}

// Transactional accepts a transaction, performs an operation (presumably on
// the active transaction), and returns an error if anything want wrong.
type Transactional func(*sqlx.Tx) error

// CreateSchema executes the required SQL for building a new schema.
// ifMissing adds the "IF NOT EXISTS" clause.
func CreateSchema(db *sqlx.DB, schema string, ifMissing bool) error {
	maybe := " "
	if ifMissing {
		maybe = " IF NOT EXISTS "
	}
	sql := fmt.Sprintf("CREATE SCHEMA%s%s", maybe, schema)
	_, err := db.Exec(sql)
	return err
}

// CreateTable builds and executes a CREATE TABLE SQL statement.
func CreateTable(db *sqlx.DB, tablename, table string, ifMissing bool) error {
	maybe := " "
	if ifMissing {
		maybe = " IF NOT EXISTS "
	}
	sql := fmt.Sprintf("Create TABLE%s%s ( %s )", maybe, tablename, table)
	_, err := db.Exec(sql)
	return err
}

// Transact wraps a DB operation in a transaction
func Transact(db *sqlx.DB, fn Transactional) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		// Return any error we had rolling back
		if txErr := tx.Rollback(); txErr != nil {
			return txErr
		}
		// Otherwise, raise the initial error so our caller knows something
		// went wrong
		return err
	}
	return tx.Commit()
}
