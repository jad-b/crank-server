package metrics

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/flagit"
	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

const (
	bodyweightTableName = "bodyweight"
	// BodyweightSQL is the SQL required to create the Bodyweight table.
	bodyweightTable = `
  user_id integer NOT NULL,
  timestamp timestamp with time zone NOT NULL,
  weight numeric(5,2) NOT NULL CHECK (weight < 1000),
  comment text,
  UNIQUE(user_id, timestamp)
`
)

// CreateTableBodyweight creates the Bodyweight table
func CreateTableBodyweight(db *sqlx.DB) error {
	return torque.CreateTable(
		db,
		Schema,
		bodyweightTableName,
		bodyweightTable,
		true)
}

// Bodyweight is a timestamped bodyweight record, with optional comment.
type Bodyweight struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Weight    float64   `json:"weight"`
	Comment   string    `json:"comment"`
}

/*
	CommandLineActor
*/

// ParseFlags parses command-line flags related to Bodyweight and loads them
// into itself.
func (bw *Bodyweight) ParseFlags(action string, args []string) {
	// Define sub-flags for the bodyweight resource
	var tsFlag flagit.TimeFlag
	bwFlags := flag.NewFlagSet("bwFlags", flag.ContinueOnError)
	bwFlags.Var(&tsFlag, "timestamp", "")

	bwFlags.Float64Var(&bw.Weight, "weight", 0.0, "")
	bwFlags.StringVar(&bw.Comment, "comment", "", "")

	// Parse the given flags
	bwFlags.Parse(args)
	// Assign our leftover timestamp
	bw.Timestamp = time.Time(tsFlag)
}

/*
	DBResource
*/

// Create inserts a new bodyweight entry into the DB.
func (bw *Bodyweight) Create(db *sqlx.DB) error {
	_, err := db.NamedExec(fmt.Sprintf(`
	INSERT INTO %s.%s (
		user_id,
		timestamp,
		weight,
		comment
	) VALUES (
		:user_id,
		:timestamp,
		:weight,
		:comment
	)`, Schema, bodyweightTableName), bw)
	return err
}

// Retrieve does a lookup for the corresponding bodyweight record by timestamp.
func (bw *Bodyweight) Retrieve(db *sqlx.DB) error {
	return db.Get(
		bw,
		fmt.Sprintf(`
		SELECT
			user_id,
			timestamp,
			weight,
			comment
		FROM %s.%s
		WHERE timestamp=$1`,
			Schema, bodyweightTableName),
		bw.Timestamp)
}

// Update modifies the matching row in the DB by timestamp.
func (bw *Bodyweight) Update(db *sqlx.DB) error {
	_, err := db.NamedExec(
		fmt.Sprintf(`
			UPDATE %s.%s
			SET
				user_id=:user_id,
				weight=:weight,
				comment=:comment
			WHERE timestamp=:timestamp`,
			Schema, bodyweightTableName),
		bw)
	return err
}

// Delete removes the row from the DB
func (bw *Bodyweight) Delete(db *sqlx.DB) error {
	stmt := fmt.Sprintf(`
			DELETE FROM %s.%s
			WHERE timestamp=:timestamp`,
		Schema, bodyweightTableName)
	_, err := db.NamedExec(stmt, bw)
	return err
}

/*
	RESTfulHandler
*/

// HandlePost creates a new bodyweight record.
func (bw *Bodyweight) HandlePost(w http.ResponseWriter, req *http.Request) {
	log.Print("Request: Create Bodyweight")
	err := torque.ReadJSONRequest(req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = bw.Create(torque.DB); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Created %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// HandleGet returns the related bodyweight record
func (bw *Bodyweight) HandleGet(w http.ResponseWriter, req *http.Request) {
	log.Print("Request: Retrieve Bodyweight")
	torque.LogRequest(req)
	var err error
	bw.Timestamp, err = torque.GetTimestampQuery(req)
	if err != nil {
		log.Print(err)
		http.Error(w, "Missing timestamp in query parameters", http.StatusBadRequest)
		return
	}
	log.Printf("Retrieving %+v", bw)
	if err := bw.Retrieve(torque.DB); err != nil {
		torque.BadRequest(w, req, "No record found")
		return
	}
	log.Printf("Retrieved %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// HandlePut updates a Bodyweight resource.
func (bw *Bodyweight) HandlePut(w http.ResponseWriter, req *http.Request) {
	// Parse body of PUT request into a Bodyweight struct
	err := torque.ReadJSONRequest(req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = bw.Update(torque.DB); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Updated %+v", bw)
	// Write updated record to client
	torque.WriteOkayJSON(w, bw)
}

// HandleDelete removes the bodyweight record from the database.
func (bw *Bodyweight) HandleDelete(w http.ResponseWriter, req *http.Request) {
	// Retrieve timestamp from request
	timestamp, err := torque.GetTimestampQuery(req)
	if err != nil {
		http.Error(w, "Invalid timestamp provided", http.StatusBadRequest)
		return
	}
	if err = bw.Delete(torque.DB); err != nil {
		http.NotFound(w, req)
		return
	}
	log.Printf("Deleted bodyweight @ %s", timestamp)
	torque.WriteOkayJSON(w, bw)
}

/*
	RESTfulResource
*/

// GetResourceName returns the name the resource wishes to be refered to by in
// the URL
func (bw *Bodyweight) GetResourceName() string {
	return torque.SlashJoin(Category, "bodyweight/")
}

// RegisterURL sets up the handler for the Bodyweight reosurce on the server.
func (bw *Bodyweight) RegisterURL() error { return nil }
