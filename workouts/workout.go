package workouts

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	workoutTableSQL = `
workout_id serial PRIMARY KEY,
user_id int REFERENCES users.user_auth (id) ON UPDATE CASCADE ON DELETE CASCADE,
last_modified timestamptz default now()
`
)

var workoutTableName = fmt.Sprintf("%s.workout", Schema)

// Create inserts a new workout row into the DB.
func (w *Workout) Create(tx *sqlx.Tx) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			user_id,
			last_modified
		) VALUES (
			$1,
			$2
		) RETURNING workout_id`, workoutTableName)
	// Set workout ID from assigned row ID
	var rowInt int64
	err := tx.QueryRowx(q, w.UserID, w.LastModified).Scan(&rowInt)
	if err != nil {
		return err
	}
	w.ID = int(rowInt)
	// Insert Exercises
	for _, ex := range w.Exercises {
		ex.WorkoutID = w.ID
		if err := ex.Create(tx); err != nil {
			return err
		}
	}
	// TODO link tags
	return err
}

// Retrieve queries the DB for a workout by ID or UserID & timestamp
func (w *Workout) Retrieve(tx *sqlx.Tx) error {
	tmp := *w
	q := fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE %s`,
		workoutTableName, w.buildWhere())
	if err := tx.Get(w, q, &tmp); err != nil {
		return err
	}
	// Retrieve associated exercise IDs
	var exIDs []int
	q := fmt.Sprintf(`
		SELECT exercise_id
		FROM %s
		WHERE workout_id=$1`,
		exerciseTableName, w.ID)
	if err := tx.Select(&exIDs, q); err != nil {
		return err
	}
	// Lookup exercises by the ID
	for _, i := range exIDs {
		ex := Exercise{ID: i}
		ex.Retrieve(tx) // Implicit set lookup
		w.Exercises = append(w.Exercises, ex)
	}
	// TODO Lookup tags
}

// Delete removes the workout from the table. It also removes any tags it held
// from the linking table.
func (w *Workout) Delete(tx *sqlx.Tx) error {
	q := fmt.Sprintf(`
		DELETE FROM %s
		WHERE workout_id=:workout_id AND user_id=:user_id
		`, workoutTableName)
	res, err := tx.NamedExec(q, w)
	// DELETE had no affect
	if i, _ := res.RowsAffected(); i == 0 {
		return errors.New("Resource does not exist")
	}
	// TODO remove linked tags
	return err
}
