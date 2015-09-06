package torque

import "github.com/jmoiron/sqlx"

// A collection of test utilities.

// Connect sets up a DB connection for the test.
// Only tests should need to frequently setup DB connections.
// This also makes the assumption you want to use the '-psql-conf' flag.
//
// Don't forgot to `defer db.Close()`
func Connect() *sqlx.DB {
	// Setup our database connection
	pgConf := LoadPostgresConfig(*PsqlConf)
	return OpenDBConnection(pgConf)
}
