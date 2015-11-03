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
			workout_id,
			user_id,
			last_modified
		) VALUES (
			:workout_id,
			:user_id,
			:last_modified
		)`, workoutTableName)
	for _, ex := range w.Exercises {
		ex.WorkoutID = w.ID
		if err := ex.Create(tx); err != nil {
			return err
		}
	}
	// TODO link tags
	_, err := tx.NamedExec(q, w)
	return err
}

// Delete removes the workout from the table. It also removes any tags it held
// from the linking table.
func (w *Workout) Delete(tx *sqlx.Tx) error {
	q := fmt.Sprintf(`
		DELETE FROM %s
		WHERE
			workout_id=:workout_id,
			user_id=:user_id
		`, workoutTableName)
	res, err := tx.NamedExec(q, w)
	// DELETE had no affect
	if i, _ := res.RowsAffected(); i == 0 {
		return errors.New("Resource does not exist")
	}
	// TODO remove linked tags
	return err
}
