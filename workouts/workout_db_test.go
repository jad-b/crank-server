package workouts

import (
	"testing"

	"github.com/jad-b/torque"
)

func TestCreateTableWorkout(t *testing.T) {
	torque.CreateSchema(db, Schema, true)

	err := torque.CreateTable(
		db,
		workoutTableName,
		workoutTableSQL,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkoutCreate(t *testing.T) {
	workout := *testWorkout
	workout.UserID = 93 // redteam user ID; needs to be dynamic
	if err := torque.Transact(db, workout.Create); err != nil {
		t.Fatal(err)
	}
	if err := torque.Transact(db, workout.Delete); err != nil {
		t.Fatal(err)
	}
}
