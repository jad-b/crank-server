package workouts

import (
	"testing"

	"github.com/jad-b/torque"
)

// redteam user ID; needs to be dynamic
// TODO make dynamic
var testUserID = 93

func TestCreateTableWorkout(t *testing.T) {
	torque.CreateSchema(db, Schema, true)

	err := torque.CreateTable(
		db,
		workoutTableName,
		workoutTableSQL,
		true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkoutCreate(t *testing.T) {
	workout := *testWorkout
	workout.UserID = testUserID
	if err := torque.Transact(db, workout.Create); err != nil {
		t.Fatal(err)
	}
	if err := torque.Transact(db, workout.Delete); err != nil {
		t.Fatal(err)
	}
}

func TestWorkoutRetrieve(t *testing.T) {
	workout := *testWorkout
	workout.UserID = testUserID

	if err := torque.Transact(db, workout.Create); err != nil {
		t.Fatal(err)
	}

	// Try to retrieve a workout
	t.Logf("Created workout = %d", workout.ID)
	wkt2 := Workout{ID: workout.ID}
	t.Logf("Attempting to retrieve workout %d", wkt2.ID)
	if err := torque.Transact(db, wkt2.Retrieve); err != nil {
		t.Fatal(err)
	}

	// Verify fields were updated
	if wkt2.Exercises == nil {
		t.Fatal("Failed to retrieve exercises")
	}
	for i := range wkt2.Exercises {
		if wkt2.Exercises[i].Sets == nil {
			t.Logf("Exercise:\n%#v", wkt2.Exercises[i])
			t.Fatalf("Failed to retrieve exercise %d sets", i)
		}
	}
	// For lack of a better quick comparison...
	t.Logf("%#v", wkt2)

	if err := torque.Transact(db, workout.Delete); err != nil {
		t.Fatal(err)
	}
}
