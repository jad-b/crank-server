package api

import (
	"net/http"
	"time"

	"github.com/jad-b/torque"
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

/*
MUSCLE_GROUPS = (('Q', 'Quadriceps'),
                 ('LB', 'Lower Back'),
                 ('B', 'Biceps'),
                 ('C', 'Chest'),
                 ('A', 'Abdominals'),
                 ('H', 'Hamstrings'),
                 ('T', 'Triceps'),
                 ('TR', 'Traps'),
                 ('M', 'Middle Back'),
                 ('L', 'Lats'),
                 ('N', 'Neck'),
                 ('F', 'Forearms'),
                 ('G', 'Glutes'),
                 ('S', 'Shoulders'),
                 ('CA', 'Calves'))
*/

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

	// Primary muscle group impacted by this exercise
	PrimaryMuscleGroup string `json:"primary_muscle_group"`

	// Secondary muscle group impacted by this exercise
	SecondaryMuscleGroup string `json:"secondary_muscle_group"`

	// What "type" of exercise is this?
	Category string `json:"category"`

	// Aliases for this Exercise
	AlsoKnownAs string `json:"also_known_as"`

	// One of Compound or Isolation, describing the mechanics of this exercise
	MechanicsType string `json:"mechanics_type"`

	// A link to a heatmap image of the muscles impacted by this exercise
	MuscleImage string `json:"muscle_image"`

	// A newline separated textual guide to this exercise
	Guide string `json:"guide"`

	// One of Beginner, Intermediate, Expert. Indicating the level of skill
	// required to perform this exercise successfully
	DifficultyLevel string `json:"difficulty_level"`

	// Free form text describing any equipment required to perform this
	// exercise
	Equipment string `json:"equipment"`

	// The direction of force applied to perform this exercise. Will be one of
	// Push, Pull, Static, or None
	Force string `json:"force"`

	// This Exercise's type
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
