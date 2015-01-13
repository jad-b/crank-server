package torque

import "github.com/jmoiron/sqlx"

// Context wraps configuration and services inside Torque. It is intended to be
// shared amongst goroutines/threads, and items provided inside the Context
// should be made concurrency-safe.
type Context struct {
	// Database connection
	DB *sqlx.DB
}
