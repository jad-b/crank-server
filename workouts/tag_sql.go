package workouts

import "github.com/jmoiron/sqlx"

var (
	// TagLinkingTableSQL describes a many-to-many table that links arbitrary
	// tag data to other rows via their row ID.
	// Opting to use an 'int' instead of a foreign key for 'tagged_id' lets us use
	// this table to link to any other table only on row id, *but* also gives up
	// the ability to do cascading updates. This may prove to be too painful to
	// keep up, and the table might have to get split.
	TagLinkingTableSQL = `
	link_id serial PRIMARY KEY,
	tagged_id int NOT NULL,
	tag_id REFERENCES tag (tag_id) ON UPDATE CASCADE,
`
	// TagTableSQL provides a basic key=value storage place.
	TagTableSQL = `
	tag_id serial PRIMARY KEY,
	tag text NOT NULL,
	value text NOT NULL
`
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

// CleanupTags removes any tags links for the tagged ID.
func CleanupTags(tx *sqlx.DB, taggedID int) error {
	m := map[string]interface{}{"tagged_id": taggedID}
	_, err := tx.NamedExec(`
		DELETE FROM workout.tag_link
		WHERE tagged_id = :tagged_id`, m)
	return err
}
