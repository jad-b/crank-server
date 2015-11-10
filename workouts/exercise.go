package workouts

import (
	"errors"
	"fmt"
	"log"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

var (
	exerciseTableName = fmt.Sprintf("%s.exercise", Schema)
	exerciseTableSQL  = `
exercise_id serial PRIMARY KEY,
workout_id int REFERENCES workouts.workout (workout_id) ON UPDATE CASCADE ON DELETE CASCADE,
movement text,
last_modified timestamptz default now()
`
	// exerciseModifiersTableSQL links exercise instances to a list of
	// available modifiers.
	exerciseModifiersTableName = torque.VariadicJoin(".", Schema, "exercise_modifiers")
	exerciseModifiersTableSQL  = `
modifier_id serial PRIMARY KEY,
modifier text
`
	exerciseModifiersLinkingTableSQL = `
exercise_id int REFERENCES workouts.exercise (exercise_id) ON UPDATE CASCADE,
modifier_id int REFERENCES workouts.exercise_modifiers (modifier_id) ON UPDATE CASCADE ON DELETE CASCADE,
CONSTRAINT exercise_modifier_key PRIMARY KEY (exercise_id, modifier_id)
`
)

// Create adds an Exercise row into the transaction
func (ex *Exercise) Create(tx *sqlx.Tx) error {
	// Insert exercise row
	q := fmt.Sprintf(`
		INSERT INTO %s (
			movement,
			last_modified
		) VALUES (
			$1,
			$2
		) RETURNING exercise_id`, exerciseTableName)
	var rowInt int64
	log.Print("Creating exercise record in DB")
	err := tx.QueryRowx(q, ex.Movement, ex.LastModified).Scan(&rowInt)
	if err != nil {
		return err
	}
	ex.ID = int(rowInt) // Downcast from int64 to int
	log.Printf("Created row %d in Exercise table", ex.ID)
	// Update Modifiers table
	err = linkModifiers(tx, ex.ID, ex.Modifiers)
	if err != nil {
		return err
	}
	// Update tags table
	err = linkTags(tx, ex.ID, ex.Tags)
	// Create Set entries
	for _, set := range ex.Sets {
		set.ExerciseID = ex.ID
		err = set.Create(tx)
		if err != nil {
			return err
		}
	}
	return err
}

// Delete removes the exercise entry.
func (ex *Exercise) Delete(tx *sqlx.Tx) error {
	q := fmt.Sprintf(`
		DELETE FROM %s
		WHERE exercise_id=:exercise_id
	`, exerciseTableName)
	res, err := tx.NamedExec(q, ex) // DELETE had no affect
	if i, _ := res.RowsAffected(); i == 0 {
		return errors.New("Resource does not exist")
	}
	// Remove modifier linkings
	// Remove tag linkings
	return err
}
