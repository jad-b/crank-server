package workouts

import "time"

// Defines the data model for workouts

// MassUnit defines units of measurement for mass
type MassUnit int

// See MassUnit
const (
	Pounds MassUnit = iota
	Kilograms
	Lbs = Pounds
	Kgs = Kilograms
)

// RepUnit is the unit of the work quantity performed.
// A weight lifter would perform repetitions, a swimmer may perform laps,
// a stretch may be for an amount of time, etc.
type RepUnit string

// See RepUnit
const (
	Repetition = "repetition"
	Second     = "second"
)

// Workout defines an entire workout
type Workout struct {
	// Owner of workout
	UserID int `json:"user_id" db:"user_id"`
	// Time of last modification
	LastModified time.Time `json:"last_modified" db:"last_modified"`
	// Exercises performed during the workout
	// one2many relationship
	Exercises []Exercise
	// Arbitrary key=value data
	// many2many relationship
	Tags []Tag
}

// Exercise is a performed (or planned) instance of an exercise
type Exercise struct {
	// Instance ID
	ID int `json:"exercise_id" db:"exercise_id"`
	// Name of the primary movement, e.g. Squat
	Name string `json:"exercise_name" db:"exercise_name"`
	// Modifiers to the movement. For Squat, you'd have Front, Box, Partial,
	// Anderson, etc.
	// m2m relationship
	Modifiers []string
	// Sets performed for the exercise
	// one2many relationship
	Sets []Set
	// Arbitrary key=value data
	// m2m relationship
	Tags         []Tag
	LastModified time.Time `json:"last_modified" db:"last_modified"`
}

// Set is a performed (or planned) workout set of an Exercise
type Set struct {
	SetID int `json:"set_id" db:"set_id"`
	// Parent exercise instance
	ExerciseID int `json:"exercise_id" db:"exercise_id"`
	Weight     int
	// Pounds, kilograms, stone, what-have-you
	WeightUnit MassUnit `json:"weight_unit" db:"weight_unit"`
	Reps       int
	RepUnit    RepUnit `json:"rep_unit" db:"rep_unit"`
	// Rest period taken *before* this set. Knowing the rest period taken
	// *after* this set would have no meaning, although it's still pretty empty
	// when taken alone.
	// A negative duration indicates the time is unknown
	Rest time.Duration
	// Number marking the order the set was performed within the workout
	// Thus, it only has meaning with the context of its parent workout
	Order int
}

// Tag holds arbitrary key=value strings
type Tag struct {
	Name  string
	Value string
}
