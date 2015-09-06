package users

import (
	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

// package.go contains package-level conts, vars, and functions.

// Schema is the DB schema name for the users package
const Schema = "users"

// CreateSchema issues the schema creation SQL statement for the users package.
func CreateSchema(db *sqlx.DB) error {
	err := torque.CreateSchema(db, Schema, true)
	return err
}
