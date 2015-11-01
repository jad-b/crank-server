package workouts

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	setTableName = fmt.Sprintf("%s.set", Schema)
	setTableSQL  = `
	set_id serial PRIMARY KEY,
	exercise_id int REFERENCES exercise (exercise_id) ON UPDATE CASCADE ON DELETE CASCADE,
	weight_unit text,
	reps integer,
	rep_unit text,
	rest interval MINUTE TO SECOND 0,
	order integer
`
)

// Create adds a set row insertion into the transaction
func (s *Set) Create(tx *sqlx.Tx) error {
	// Create Set row
	q := fmt.Sprintf(`
		INSERT INTO %s.%s (
			exercise_id,
			weight,
			weight_unit,
			reps,
			rep_unit,
			rest,
			order
		) VALUES (
			:exercise_id,
			:weight,
			:weight_unit,
			:reps,
			:rep_unit,
			:rest,
			:order
		)`, Schema, setTableName)
	_, err := tx.NamedExec(q, s)
	return err
}

// Retrieve looks up an exercise set row using whatever it was given.
func (s *Set) Retrieve(tx *sqlx.Tx) error {
	// Copy, so as to avoid any overwriting weirdness
	clone := *s
	q := fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE %s
	`, setTableName, s.buildWhere())
	return tx.Get(s, q, &clone)
}

// Update does just that.
func (s *Set) Update(tx *sqlx.Tx) error {
	q := fmt.Sprintf(`
		UPDATE %s
		SET
			set_id=:set_id,
			exercise_id=:exercise_id,
			weight=:weight,
			weight_unit=:weight_unit,
			reps=:reps,
			rep_unit=:rep_unit,
			rest=:rest,
			order=:order
		WHERE %s
	`, setTableName, s.buildWhere())
	_, err := tx.NamedExec(q, s)
	return err
}

// Delete removes the exercise set from the database
func (s *Set) Delete(tx *sqlx.Tx) error {
	q := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s
	`, setTableName, s.buildWhere())
	res, err := tx.NamedExec(q, s)
	// DELETE had no affect
	if i, _ := res.RowsAffected(); i == 0 {
		return errors.New("Resource does not exist")
	}
	return err
}

// Build WHERE clause based off of provided data
func (s *Set) buildWhere() string {
	if s.SetID != 0 { // Use id lookup
		return "set_id=:set_id"
	}
	return `exercise_id=:exercise_id, order=:order`
}
