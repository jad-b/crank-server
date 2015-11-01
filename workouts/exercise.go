package workouts

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	exerciseTableName = fmt.Sprintf("%s.exercise", Schema)
	// ExerciseModifiersTableSQL links exercise instances to a list of
	// available modifiers.
	ExerciseModifiersTableSQL = ``
)

// Create adds an Exercise row into the transaction
func (ex *Exercise) Create(tx *sqlx.Tx) error {
	// Insert exercise row
	q := fmt.Sprintf(`
		INSERT INTO %s.%s (
			exercise_name,
			last_modified
		) VALUES (
			:exercise_name,
			:last_modified
		)`, Schema, exerciseTableName)
	res, err := tx.NamedExec(q, ex)
	if err != nil {
		return err
	}
	// Get our row ID
	rowInt, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ex.ID = int(rowInt) // Downcast from int64 to int
	// Update Modifiers table
	err = linkModifiers(tx, ex.ID, ex.Modifiers)
	if err != nil {
		return err
	}
	// Update tags table
	err = linkTags(tx, ex.ID, ex.Tags)
	// Create Set entries
	for _, set := range ex.Sets {
		err = set.Create(tx)
		if err != nil {
			return err
		}
	}
	return err
}
