package ui

import "time"

/**
 * A copy of the key torque objects (Workout, Exercise) while the front-end
 * interactions are being hashed out.
 */

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
	Tags string `json:"tags"`
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
	Sets string
	// Modifiers to the movement. For Squat, you'd have Front, Box, Partial,
	// Anderson, etc.
	Modifiers string
	// Arbitrary key=value data
	// m2m relationship
	Tags string
}
