package workouts

import (
	"database/sql/driver"
	"time"
)

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

// Value converts the RepUnit into a string
func (ru RepUnit) Value() (driver.Value, error) {
	return driver.Value(string(ru)), nil
}

// See RepUnit
const (
	Repetition = "repetition"
	Second     = "second"
)

// Workout defines an entire workout
type Workout struct {
	ID int `json:"workout_id" db:"workout_id"`
	// Owner of workout
	UserID int `json:"user_id" db:"user_id"`
	// Time of last modification
	LastModified time.Time `json:"last_modified" db:"last_modified"`
	// Exercises performed during the workout
	// one2many relationship
	Exercises []Exercise `json:"exercises`
	// Arbitrary key=value data
	// many2many relationship
	Tags []Tag `json:"tags"`
}

// Exercise is a performed (or planned) instance of an exercise
type Exercise struct {
	// Instance ID
	ID int `json:"exercise_id" db:"exercise_id"`
	// Workout it belongs to
	WorkoutID int `json:"workout_id" db:"workout_id"`
	// Name of the primary movement, e.g. Squat
	Movement     string    `json:"movement" db:"movement"`
	LastModified time.Time `json:"last_modified" db:"last_modified"`
	// Sets performed for the exercise
	// one2many relationship
	Sets []Set
	// Modifiers to the movement. For Squat, you'd have Front, Box, Partial,
	// Anderson, etc.
	// m2m relationship
	Modifiers []string
	// Arbitrary key=value data
	// m2m relationship
	Tags []Tag
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
	Order int `db:"ordering"`
}

// Tag holds arbitrary key=value strings
type Tag struct {
	Name  string
	Value string
}
