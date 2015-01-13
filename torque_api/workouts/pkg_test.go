package workouts

import "testing"

func TestCreatingPkgTables(t *testing.T) {
	if err := CreateTables(); err != nil {
		t.Fatal(err)
	}
}
