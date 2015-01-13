// +build db
package workouts

import (
	"os"
	"testing"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

// Database connection for running workouts tests
var db *sqlx.DB

func TestMain(m *testing.M) {
	os.Exit(m.Run())
	db.Close()
}

func init() {
	db = torque.Connect()
}

func TestCreateExerciseSchema(t *testing.T) {
	if err := torque.CreateSchema(db, Schema, true); err != nil {
		t.Fatal(err)
	}
}
