package workouts

import (
	"testing"

	"github.com/jad-b/torque"
)

func TestCreateTableSet(t *testing.T) {
	torque.CreateSchema(db, Schema, true)

	err := torque.CreateTable(
		db,
		setTableName,
		setTableSQL,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetCreate(t *testing.T) {
	set := testSets[0]
	if err := torque.Transact(db, set.Create); err != nil {
		t.Fatal(err)
	}
	if err := torque.Transact(db, set.Delete); err != nil {
		t.Fatal(err)
	}
}

func TestSetRetrieveOnID(t *testing.T) {
	set := testSets[0]
	if err := torque.Transact(db, set.Create); err != nil {
		t.Fatal(err)
	}
	defer torque.Transact(db, set.Delete)

	set2 := Set{SetID: set.SetID}
	if err := torque.Transact(db, set2.Retrieve); err != nil {
		t.Fatal(err)
	}
	if set2.SetID != set.SetID {
		t.Fatalf("Retrieved wrong exercise set ID;\nGot: %s\nWanted: %s",
			set2.SetID, set.SetID)
	}
	if set2.Weight != set.Weight {
		t.Fatalf("Retrieved wrong exercise set weight;\nGot: %s\nWanted: %s",
			set2.Weight, set.Weight)
	}
	if set2.ExerciseID != set.ExerciseID {
		t.Fatalf("Retrieved wrong exercise set ExerciseID;\nGot: %s\nWanted: %s",
			set2.ExerciseID, set.ExerciseID)
	}
	if set2.Order != set.Order {
		t.Fatalf("Retrieved wrong exercise set Order;\nGot: %s\nWanted: %s",
			set2.Order, set.Order)
	}
}

func TestSetRetrieveOnExAndOrder(t *testing.T) {
	set := testSets[0]
	if err := torque.Transact(db, set.Create); err != nil {
		t.Fatal(err)
	}
	defer torque.Transact(db, set.Delete)

	set2 := Set{ExerciseID: set.ExerciseID, Order: set.Order}
	if err := torque.Transact(db, set2.Retrieve); err != nil {
		t.Fatal(err)
	}
	if set2.SetID != set.SetID {
		t.Fatalf("Retrieved wrong exercise set;\nGot: %s\nWanted: %s",
			set2.SetID, set.SetID)
	}
	if set2.Weight != set.Weight {
		t.Fatalf("Retrieved wrong exercise set;\nGot: %s\nWanted: %s",
			set2.Weight, set.Weight)
	}
	if set2.ExerciseID != set.ExerciseID {
		t.Fatalf("Retrieved wrong exercise set;\nGot: %s\nWanted: %s",
			set2.ExerciseID, set.ExerciseID)
	}
	if set2.Order != set.Order {
		t.Fatalf("Retrieved wrong exercise set;\nGot: %s\nWanted: %s",
			set2.Order, set.Order)
	}
}

func TestSetUpdate(t *testing.T) {
	// Create
	set := testSets[0]
	if err := torque.Transact(db, set.Create); err != nil {
		t.Fatal(err)
	}
	defer torque.Transact(db, set.Delete)

	// Update
	set1 := set // Make a copy of the shared struct
	set1.Weight = 200
	if err := torque.Transact(db, set1.Update); err != nil {
		t.Fatal(err)
	}

	// Retrieve
	set2 := Set{SetID: set.SetID}
	if err := torque.Transact(db, set.Retrieve); err != nil {
		t.Fatal(err)
	}
	if set2.Weight != set1.Weight {
		t.Fatal("Failed to update Weight")
	}
}

func TestSetDelete(t *testing.T) {
	set := testSets[0]
	// Create
	if err := torque.Transact(db, set.Create); err != nil {
		t.Fatal(err)
	}
	// Delete
	if err := torque.Transact(db, set.Delete); err != nil {
		t.Fatal(err)
	}
	// Retrieve (and fail)
	if err := torque.Transact(db, set.Retrieve); err == nil {
		t.Fatal("Failed to delete set")
	}

}
