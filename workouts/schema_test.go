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
			Exercise{}, // Placeholder for test exercise
			{
				ID:        3,
				Movement:  "Curl",
				Modifiers: []string{"Ring"},
				Sets: []Set{
					{2, 3, 177, Pounds, 10, Repetition, -1, 8},
					{3, 3, 177, Pounds, 5, Repetition, -1, 9},
				},
				LastModified: time.Now(),
			},
		},
		Tags: []Tag{},
	}
	testExercise = &Exercise{
		ID:           2,
		Movement:     "Squat",
		Modifiers:    []string{"Back"},
		Sets:         nil, // Placeholder for testSet
		Tags:         []Tag{{"comment", "coming off drill weekend, tired and small"}},
		LastModified: time.Now(),
	}
	testSets = []Set{
		{
			SetID:      4,
			ExerciseID: testExercise.ID,
			Weight:     20,
			WeightUnit: Kilograms,
			Reps:       5,
			RepUnit:    Repetition,
			Rest:       -1,
			Order:      1,
		},
		{5, testExercise.ID, 60, Kilograms, 5, Repetition, -1, 2},
		{6, testExercise.ID, 80, Kilograms, 3, Repetition, -1, 3},
		{7, testExercise.ID, 90, Kilograms, 3, Repetition, -1, 4},
		{8, testExercise.ID, 91, Kilograms, 5, Repetition, -1, 5},
		{9, testExercise.ID, 105, Kilograms, 5, Repetition, -1, 6},
		{10, testExercise.ID, 119, Kilograms, 4, Repetition, -1, 7},
	}
)

func init() {
	testWorkout.Exercises[1] = *testExercise
	testExercise.Sets = testSets
}

func TestWorkoutCreated(t *testing.T) {
	if testSets != nil && testExercise != nil && testWorkout != nil {
		t.Log("Workout compiled successfully")
	}
}
