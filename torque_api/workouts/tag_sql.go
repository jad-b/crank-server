package workouts

import (
	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

const (
	// tagTableSQL provides a basic key=value storage place.
	tagTableSQL = `
	tag_id serial PRIMARY KEY,
	tag text NOT NULL,
	value text NOT NULL`

	// workoutTagTableSQL describes a many-to-many table that links arbitrary
	// tag data to other rows via their row ID.
	workoutTagTableSQL = `
	workout_id int REFERENCES workouts.workout (workout_id) ON UPDATE CASCADE ON DELETE CASCADE,
	tag_id int REFERENCES workouts.tag (tag_id) ON UPDATE CASCADE,
	CONSTRAINT workout_tag_id PRIMARY KEY (workout_id, tag_id)`

	// exerciseTagTableSQL describes a many-to-many table that links arbitrary
	// tag data to other rows via their row ID.
	exerciseTagTableSQL = `
	exercise_id int REFERENCES workouts.exercise (exercise_id) ON UPDATE CASCADE ON DELETE CASCADE,
	tag_id int REFERENCES workouts.tag (tag_id) ON UPDATE CASCADE,
	CONSTRAINT exercise_tag_id PRIMARY KEY (exercise_id, tag_id)`
)

var (
	tagTableName         = torque.VariadicJoin(".", Schema, "tag")
	exerciseTagTableName = torque.VariadicJoin(".", Schema, "exercise_tag")
	workoutTagTableName  = torque.VariadicJoin(".", Schema, "workkout_tag")
)

// linkTags inserts a new entry into the linking table.
func linkTags(tx *sqlx.Tx, taggee int, tags []Tag) error {
	// Lookup tag IDs by name
	// Create row entries
	return nil
}

func linkModifiers(tx *sqlx.Tx, taggee int, tags []string) error {
	return nil
}
