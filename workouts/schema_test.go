package workouts

import (
	"testing"
	"time"
)

var (
	testWorkout = &Workout{
		UserID:       1,
		LastModified: time.Now(),
		Exercises: []Exercise{
			{
				ID:        1,
				Movement:  "Swing",
				Modifiers: []string{"Kettlebell"},
				Sets: []Set{
					{1, 1, 28, Kilograms, 36, Repetition, -1, 0},
				},
				LastModified: time.Now(),
			},
			{
				ID:        2,
				Movement:  "Squat",
				Modifiers: []string{"Back"},
				Sets: []Set{
					{
						SetID:      4,
						ExerciseID: 2,
						Weight:     20,
						WeightUnit: Kilograms,
						Reps:       5,
						RepUnit:    Repetition,
						Rest:       -1,
						Order:      1,
					},
					{5, 2, 60, Kilograms, 5, Repetition, -1, 2},
					{6, 2, 80, Kilograms, 3, Repetition, -1, 3},
					{7, 2, 90, Kilograms, 3, Repetition, -1, 4},
					{8, 2, 91, Kilograms, 5, Repetition, -1, 5},
					{9, 2, 105, Kilograms, 5, Repetition, -1, 6},
					{10, 2, 119, Kilograms, 4, Repetition, -1, 7},
				},
				Tags:         []Tag{{"comment", "coming off drill weekend, tired and small"}},
				LastModified: time.Now(),
			},
			{
				ID:        3,
				Movement:  "Curl",
				Modifiers: []string{"Ring"},
				Sets: []Set{
					{2, 3, 35, Pounds, 10, Repetition, -1, 8},
					{3, 3, 35, Pounds, 5, Repetition, -1, 9},
				},
				LastModified: time.Now(),
			},
		},
		Tags: []Tag{},
	}
	testExercise = &testWorkout.Exercises[1]
	testSets     = testExercise.Sets
)

func TestWorkoutCreated(t *testing.T) {
	if testSets != nil && testExercise != nil && testWorkout != nil {
		t.Log("Workout compiled successfully")
	}
}
