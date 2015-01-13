// +build test db

package metrics

import (
	"testing"
	"time"

	"github.com/jad-b/torque"
)

var (
	testBodyweight = &Bodyweight{
		UserID:    1,
		Weight:    182.4,
		Timestamp: time.Now(),
		Comment:   "This is only a test",
	}
)

func TestCreateBodyweightSchema(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	if err := CreateSchema(db); err != nil {
		t.Fatal(err)
	}
}

func TestCreateBodyweightTable(t *testing.T) {
	db := torque.Connect()
	defer db.Close()
	CreateSchema(db)

	if err := CreateTableBodyweight(db); err != nil {
		t.Fatal(err)
	}
}

func TestBodyweightCreate(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	bw := *testBodyweight
	bw.Comment = "TestBodyweightCreate"
	if err := bw.Create(db); err != nil {
		t.Fatal(err)
	}
	// Delete
	if err := bw.Delete(db); err != nil {
		t.Fatal(err)
	}
}

func TestBodyweightRetrieve(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	bw := *testBodyweight
	bw.Comment = "TestBodyweightRetrieve"
	if err := bw.Create(db); err != nil {
		t.Fatal(err)
	}
	defer bw.Delete(db)
	// Retrieve
	bw2 := Bodyweight{
		Weight:    190.4, // This value should get overridden
		Timestamp: bw.Timestamp,
	}
	if err := bw2.Retrieve(db); err != nil {
		t.Fatal(err)
	}
	// Verify it's the right bw.
	if bw2.Comment != bw.Comment {
		t.Fatalf("Retrieved wrong bw;\nGot: %s\nWanted: %s",
			bw2.Comment, bw.Comment)
	}
	if bw2.Weight != bw.Weight {
		t.Fatalf("Failed to override Weight being '%f'", bw2.Weight)
	}
	if bw2.UserID != bw.UserID {
		t.Fatal("Failed to retrieve User ID")
	}
}

func TestBodyweightUpdate(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	bw := *testBodyweight
	bw.Comment = "TestBodyweightUpdate"
	if err := bw.Create(db); err != nil {
		t.Fatal(err)
	}
	defer bw.Delete(db)

	// Update
	bw1 := bw // Make a copy of the shared struct
	bw1.Weight = 190.4
	bw1.Update(db)

	// Retrieve
	bw2 := Bodyweight{Timestamp: bw.Timestamp}
	if err := bw2.Retrieve(db); err != nil {
		t.Fatal(err)
	}
	if bw2.Weight != bw1.Weight {
		t.Fatal("Failed to update Weight")
	}
}

func TestBodyweightDelete(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	bw := *testBodyweight
	bw.Comment = "TestBodyweightDelete"
	if err := bw.Create(db); err != nil {
		t.Fatal(err)
	}
	// Delete
	if err := bw.Delete(db); err != nil {
		t.Fatal(err)
	}
	// Retrieve (and fail)
	if err := bw.Retrieve(db); err == nil {
		t.Fatal("Failed to delete bw")
	}
}
