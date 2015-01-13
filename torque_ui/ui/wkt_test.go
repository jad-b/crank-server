package ui

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func ExpectedGot(expected, got, msg string) string {
	errMsg := fmt.Sprintf("Expected v. Got\n'%s' v. '%s'", expected, got)
	if msg != "" {
		errMsg = msg + "\n" + errMsg
	}
	return errMsg
}

func TestStringToWktTime(t *testing.T) {
	timestamp := "29 Dec 15 08:26 EST"
	wktTime, err := StringToWktTime(timestamp)
	if err != nil {
		t.Fatal(err)
	}
	expect := "2015 Dec 29 @ 0826"
	if wktTime != expect {
		t.Fatalf("Expected v. Got: %s v. %s", expect, wktTime)
	}
}

func TestWorkoutToWkt(t *testing.T) {
	if len(testWorkout.Exercises) != 6 {
		t.Fatal("Somehow you assigned less than 6 exercises")
	}
	wkt, err := WorkoutToWkt(testWorkout)
	if err != nil {
		t.Fatal(err)
	}
	// Split our workout according to newlines
	wktLines := strings.Split(wkt, "\n")
	expected := []string{
		// Workout Header
		then.Format(WktTimeLayout),
		"- unit: reps x kgs",
		"- this is a comment",
		// Exercises
		fmt.Sprintf("%s: %s", testWorkout.Exercises[0].Movement, testWorkout.Exercises[0].Sets),
		"- Training Max: 230 lbs",
		"- week: 3",
		fmt.Sprintf("%s: %s", testWorkout.Exercises[1].Movement, testWorkout.Exercises[1].Sets),
		"- prev: 49 x 8/5",
		fmt.Sprintf("%s: %s", testWorkout.Exercises[2].Movement, testWorkout.Exercises[2].Sets),
		"- prev: 49 x 7/4",
		fmt.Sprintf("%s: %s", testWorkout.Exercises[3].Movement, testWorkout.Exercises[3].Sets),
		"- prev: 7 x 8/4",
		fmt.Sprintf("%s: %s", testWorkout.Exercises[4].Movement, testWorkout.Exercises[4].Sets),
		"- prev: 20 x 7/5",
		fmt.Sprintf("%s: %s", testWorkout.Exercises[5].Movement, testWorkout.Exercises[5].Sets),
		"- prev: 35 x 10/4",
		"- unit: lbs x reps",
	}
	for i, exp := range expected {
		got := wktLines[i]
		t.Logf("'%s' =?= '%s'", exp, got)
		if exp != got {
			t.Error(ExpectedGot(exp, got, ""))
		}
	}
}

func TestMinimalWorkoutToWkt(t *testing.T) {
	var w Workout
	if err := json.Unmarshal([]byte(emptyWorkout), &w); err != nil {
		t.Fatal(err)
	}
	wkt, err := WorkoutToWkt(&w)
	if err != nil {
		t.Fatal(err)
	}
	badOutput := "2015 Dec 31 @ 0058\n- \n: \n- \n"
	if wkt == badOutput {
		t.Error("Still getting unwanted output: " + wkt)
	} else {
		t.Log("Wkt: " + wkt)
	}
}
