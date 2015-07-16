package api

/*
	URL Path: /workout
*/

import (
	"net/http"
	"time"

	"github.com/jad-b/torque"
)

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

// GetWorkout returns a workout by timestamp
func GetWorkout(w http.ResponseWriter, req *http.Request) {
	timestamp, err := web.Stamp(req)
	workout, err := crank.LookupWorkout(timestamp)
	if err != nil {
		http.NotFound(w, req) // Write 404 to response
		return
	}
	writeJSON(w, http.StatusOK, workout)
}
