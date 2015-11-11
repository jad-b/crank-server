package workouts

import (
	"testing"

	"github.com/jad-b/torque"
)

func TestCreateTableExercise(t *testing.T) {
	torque.CreateSchema(db, Schema, true)

	err := torque.CreateTable(
		db,
		exerciseTableName,
		exerciseTableSQL,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExerciseCreate(t *testing.T) {
	ex := *testExercise
	if err := torque.Transact(db, ex.Create); err != nil {
		t.Fatal(err)
	}
	if err := torque.Transact(db, ex.Delete); err != nil {
		t.Fatal(err)
	}
}

func TestSetRetrieval(t *testing.T) {
	ex := *testExercise
	if err := torque.Transact(db, ex.Create); err != nil {
		t.Fatal(err)
	}
	if err := torque.Transact(db, ex.Delete); err != nil {
		t.Fatal(err)
	}
}
