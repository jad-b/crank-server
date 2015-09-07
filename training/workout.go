package api

import (
	"net/http"
	"time"
)

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

// Collection of valid values for an assortment of ExerciseMeta fields. Note
// that these can not be marked as constant, because Array's are mutable.
// Because of this avoid, at all costs, ever writing to any of these Arrays
var (
	ExerciseCategories = [3]string{"Resistance", "Functional", "Cardio"}
	MuscleGroups       = [15]string{"Quadriceps", "Lower Back", "Biceps", "Chest",
		"Abdominals", "Hamstrings", "Triceps", "Traps",
		"Middle Back", "Lats", "Neck", "Forearms",
		"Glutes", "Shoulders", "Calves"}
	ExerciseMechanics = [2]string{"Compound", "Isolation"}
	DifficultyLevels  = [3]string{"Beginner", "Intermediate", "Expert"}
	Forces            = [3]string{"Push", "Pull", "Static"}
	ExerciseTypes     = [7]string{"Cardio", "Olympic Weightlifting", "Plyometrics",
		"Powerlifting", "Strength", "Stretching",
		"Strongman"}
)

// ExerciseMeta is used to encapsulate the metadata of a specific exercise.
type ExerciseMeta struct {
	// The name of this exercise
	Name string `json:"name"`

	// The rate at which weight/reps for this excercise should be increased by
	ProgressDifferential float64 `json:"progress_differential"`

	// Primary muscle group impacted by this exercise. See MuscleGroups for
	// the list of valid muscle group values
	PrimaryMuscleGroup string `json:"primary_muscle_group"`

	// Secondary muscle group impacted by this exercise. See MuscleGroups for
	// the list of valid muscle group values
	SecondaryMuscleGroup string `json:"secondary_muscle_group"`

	// What "type" of exercise is this? See ExerciseCategories for list of
	// valid values
	Category string `json:"category"`

	// Aliases for this Exercise
	AlsoKnownAs string `json:"also_known_as"`

	// Describes the mechanics of this exercise. See ExerciseMechanics for list
	// of valid values
	MechanicsType string `json:"mechanics_type"`

	// A link to a heatmap image of the muscles impacted by this exercise
	MuscleImage string `json:"muscle_image"`

	// A newline separated textual guide to this exercise
	Guide string `json:"guide"`

	// Indicates the level of skill, on average, required to perform this
	// exercise successfully. See DifficultyLevels for list of valid values
	DifficultyLevel string `json:"difficulty_level"`

	// Free form text describing any equipment required to perform this
	// exercise
	Equipment string `json:"equipment"`

	// The direction of force applied to perform this exercise. See Forces for
	// the list of valid values
	Force string `json:"force"`

	// This Exercise's type. See ExerciseTypes for valid values
	ExerciseType string `json:"type"`
}

// Set ...
type Set struct {
	Reps   uint8 `json:"reps"`
	Weight uint8 `json:"weight"`
	Order  uint8 `json:"order"`
	Rest   uint8 `json:"rest"`
}

// Get returns a workout by timestamp
func Get(w http.ResponseWriter, req *http.Request) {
	timestamp, err := web.Stamp(req)
	workout, err := crank.LookupWorkout(timestamp)
	if err != nil {
		http.NotFound(w, req) // Write 404 to response
		return
	}
	writeJSON(w, http.StatusOK, workout)
}
