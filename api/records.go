package api

import "time"

/*
	Metrics
*/

// BodyweightRecord is a timestamped bodyweight record.
type BodyweightRecord struct {
	Bodyweight float32 `json:"bodyweight"`
	// 'omitempty' => Skip the timestamp if it's empty
	Timestamp time.Time `json:"timestamp"`
}

/*
	Workouts
*/

// Workout ...
type Workout struct {
	Timestamp time.Time  `json:"timestamp"`
	Comment   string     `json:"comment"`
	Exercises []Exercise `json:"exercises"`
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
