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

// Workout defines an entire workout
type Workout struct {
	// Owner of workout
	UserID int
	// Time of last modification
	LastModified time.Time
	// Exercises performed during the workout
	Exercises []Exercise
	// Arbitrary key=value data
	Tags []Tag
}

// Exercise is a performed (or planned) instance of an exercise
type Exercise struct {
	// Name of the primary movement, e.g. Squat
	Name string
	// Modifiers to the movement
	// for Squat you'd have Front, Box, Partial, Anderson, etc.
	Modifiers []string
	// Sets performed for the exercise
	Sets []Set
	// Arbitrary key=value data
	Tags         []Tag
	LastModified time.Time
}

// Set is a performed (or planned) workout set of an Exercise
type Set struct {
	Weight int
	// Pounds, kilograms, stone, what-have-you
	WeightUnit MassUnit
	Reps       int
	// Rest period taken *before* this set. Knowing the rest period taken
	// *after* this set would have no meaning, although it's still pretty empty
	// when taken alone.
	// A negative duration indicates the time is unknown
	Rest time.Duration
	// Number marking the order the set was performed within the workout
	// Thus, it only has meaning with the context of its parent workout
	SetID int
}

// Tag holds arbitrary key=value strings
type Tag struct {
	Name  string
	Value string
}
