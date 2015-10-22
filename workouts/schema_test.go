package workouts

import (
	"testing"
	"time"
)

var (
	s = &[]Set{
		{20, Kilograms, 5, -1, 1},
		{60, Kilograms, 5, -1, 2},
		{80, Kilograms, 3, -1, 3},
		{90, Kilograms, 3, -1, 4},
		{91, Kilograms, 5, -1, 5},
		{105, Kilograms, 5, -1, 6},
		{119, Kilograms, 4, -1, 7},
	}
	ex = &Exercise{
		Name:         "Squat",
		Modifiers:    []string{"Back"},
		Sets:         *s,
		Tags:         []Tag{{"comment", "coming off drill weekend, tired and small"}},
		LastModified: time.Now(),
	}
	wkt = &Workout{
		UserID:       1,
		LastModified: time.Now(),
		Exercises: []Exercise{
			{
				Name:      "Swing",
				Modifiers: []string{"Kettlebell"},
				Sets: []Set{
					{28, Kilograms, 36, -1, 0},
				},
				LastModified: time.Now(),
			},
			*ex,
			{
				Name:      "Curl",
				Modifiers: []string{"Ring"},
				Sets: []Set{
					{177, Pounds, 10, -1, 8},
					{177, Pounds, 5, -1, 9},
				},
				LastModified: time.Now(),
			},
		},
		Tags: []Tag{},
	}
)

func TestWorkoutCreated(t *testing.T) {
	if s != nil && ex != nil && wkt != nil {
		t.Log("Workout built successfully")
	}
}
