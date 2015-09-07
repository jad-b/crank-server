// +build test db

package redteam

import (
	"fmt"
	"testing"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

/*
	sql_test runs a set of basic DB operations to ensure things work as you might
expect. Normally I wouldn't test library code, but this is as much a learning
experience with Go & SQL as it is a guide to future developers.
*/

var testSchema = "testing"

func createTestingSchema(db *sqlx.DB, schema string) error {
	createSchemaSQL := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)
	_, err := db.Exec(createSchemaSQL)
	return err
}

func TestCreateTestingSchema(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	err := createTestingSchema(db, testSchema)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTable(t *testing.T) {
	db := torque.Connect()
	defer db.Close()
	// Setup req'd schema; ignore errors
	createTestingSchema(db, testSchema)

	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s.person (
			first_name text,
			last_name text,
			email text
		);

		CREATE TABLE IF NOT EXISTS %[1]s.place (
			country text,
			city text NULL,
			telcode integer
		)`, testSchema)
	_, err := db.Exec(createTableSQL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropTestingSchema(t *testing.T) {
	db := torque.Connect()
	defer db.Close()
	// Setup req'd schema; ignore errors
	createTestingSchema(db, testSchema)

	dropSchemaSQL := fmt.Sprintf("DROP SCHEMA IF EXISTS %[1]s CASCADE", testSchema)
	_, err := db.Exec(dropSchemaSQL)
	if err != nil {
		t.Fatal(err)
	}
}
