package crank

import (
	"time"
	"fmt"
)

// Workout ...
type Workout struct {
	Timestamp time.Time  `json:"timestamp"`
	Comment   string     `json:"comment"`
	Exercises []Exercise `json:"exercises"`
}

// NewWorkout creates a workout with the current time.
func NewWorkout() *Workout {
	return &Workout{Timestamp: time.Now()}
}

// LookupWorkout returns workouts from the database by timestamp.
func LookupWorkout(timestamp time.Time) (w *Workout, err error) {
	return &Workout{
		Timestamp: timestamp,
		Comment: fmt.Sprintf("Time is %s", timestamp.Format(time.RFC3339)),
	}, nil
}

// Exercise ...
type Exercise struct {
	Sets    []Set  `json:"sets"`
	Comment string `json:"comment"`
}

// ExerciseTag ...
type ExerciseTag struct {
	Tag     string `json:"tag"`
	Primary bool   `json:"primary"`
}

// Set ...
type Set struct {
	Reps   uint8 `json:"reps"`
	Weight uint8 `json:"weight"`
	Order  uint8 `json:"order"`
	Rest   uint8 `json:"rest"`
}
