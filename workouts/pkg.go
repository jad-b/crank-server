package workouts

import "github.com/jad-b/torque"

// DB constants
const (
	Schema = "workouts"
)

// CreateTables loads *all* tables for the package into the database.
func CreateTables() error {
	torque.CreateSchema(db, Schema, true)
	for _, tbl := range [][]string{
		{workoutTableName, workoutTableSQL},
		{exerciseTableName, exerciseTableSQL},
		{setTableName, setTableSQL},
		{tagTableName, tagTableSQL},
		{exerciseTagTableName, exerciseTagTableSQL},
		{workoutTagTableName, workoutTagTableSQL},
	} {
		err := torque.CreateTable(db, tbl[0], tbl[1], true)
		if err != nil {
			return err
		}
	}
	return nil
}
