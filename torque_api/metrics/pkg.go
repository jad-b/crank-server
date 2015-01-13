package metrics

import (
	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

// Constants
const (
	Schema   = "metrics"
	Category = "metrics"
)

// CreateSchema issues the schema creation SQL statement for the users package.
func CreateSchema(db *sqlx.DB) error {
	err := torque.CreateSchema(db, Schema, true)
	return err
}
