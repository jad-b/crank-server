package workouts

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

var (
	setTableName = fmt.Sprintf("%s.set", Schema)
	setTableSQL  = `
set_id serial PRIMARY KEY,
exercise_id int REFERENCES workouts.exercise (exercise_id) ON UPDATE CASCADE ON DELETE CASCADE,
weight integer,
weight_unit text,
reps integer,
rep_unit text,
rest bigint,
ordering integer
`
)

// Create adds a set row insertion into the transaction
func (s *Set) Create(tx *sqlx.Tx) error {
	// Create Set row
	log.Print("Creating set record in DB")
	q := fmt.Sprintf(`
		INSERT INTO %s (
			exercise_id,
			weight,
			weight_unit,
			reps,
			rep_unit,
			rest,
			ordering
		) VALUES (
			:exercise_id,
			:weight,
			:weight_unit,
			:reps,
			:rep_unit,
			:rest,
			:ordering
		)`, setTableName)
	_, err := tx.NamedExec(q, s)
	return err
}

// Retrieve looks up an exercise set row using whatever it was given.
func (s *Set) Retrieve(tx *sqlx.Tx) error {
	q := fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE %s
	`, setTableName, s.buildWhere())
	_, err := tx.NamedExec(q, s)
	return err
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
	res, err := tx.NamedExec(q, s) // DELETE had no affect
	if i, _ := res.RowsAffected(); i == 0 {
		return errors.New("Resource does not exist")
	}
	return err
}

// RetrieveSetsByExerciseID does that.
func RetrieveSetsByExerciseID(tx *sqlx.Tx, exerciseID int) (sets []Set, err error) {
	if exerciseID == 0 {
		return nil, errors.New("Exercise has no ID for lookup")
	}
	// Build query
	var rows *sqlx.Rows
	q := fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE exercise_id=$1`,
		setTableName)
	// Retrieve all sets for this exercise
	rows, err = tx.Queryx(q, exerciseID)
	if err != nil {
		return nil, err
	}
	// Scan them into Set structs
	i := 1
	for rows.Next() {
		var set Set
		if err = rows.StructScan(&set); err == nil {
			log.Printf("Set: %#v", set)
			sets = append(sets, set)
		} else {
			log.Print(err)
		}
		log.Printf("Scanned %d set(s)", i)
		i++
	}
	if err = rows.Err(); err != nil {
		log.Print(err)
	}
	return sets, err
}

// Build WHERE clause based off of provided data
func (s *Set) buildWhere() string {
	if s.SetID != 0 { // Use id lookup
		return "set_id=:set_id"
	}
	return `exercise_id=:exercise_id, order=:order`
}
