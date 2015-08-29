package metrics

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/users"
)

const (
	// BodyweightSQL is the SQL required to create the Bodyweight table.
	BodyweightSQL = `
CREATE TABLE metrics.bodyweight (
  "timestamp" timestamp(0) with time zone NOT NULL UNIQUE,
  weight numeric(5,2) NOT NULL CHECK (weight < 1000),
  comment text
);
`
)

// Bodyweight is a timestamped bodyweight record, with optional comment.
type Bodyweight struct {
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
	var tsFlag torque.TimestampFlag
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
func (bw *Bodyweight) Create(conn *sql.DB) error {
	_, err := conn.Exec(`
	INSERT INTO metrics.bodyweight (timestamp, weight, comment)
	VALUES ($1, $2, $3)`,
		bw.Timestamp, bw.Weight, bw.Comment)
	if err != nil {
		return err
	}
	return nil
}

// Retrieve does a lookup for the corresponding bodyweight record by timestamp.
func (bw *Bodyweight) Retrieve(conn *sql.DB) error {
	log.Printf("Looking up Bodyweight record from %s from DB", bw.Timestamp)
	err := conn.QueryRow(`
	SELECT (timestamp, weight, comment)
	FROM metrics.bodyweight
	WHERE timestamp=$1`, bw.Timestamp).Scan(bw)
	if err != nil {
		log.Printf("Problem reading from database: %s", err.Error())
		return err
	}
	return nil
}

// Update modifies the matching row in the DB by timestamp.
func (bw *Bodyweight) Update(conn *sql.DB) error {
	// Update record in database
	// TODO Only overwrite with provided fields. Maybe by building the SQL
	// statement string w/ conditional logic?
	_, err := conn.Exec(`
	UPDATE metrics.bodyweight
	SET weight=$2, comment='$3'
	WHERE timestamp > $1`,
		bw.Timestamp, bw.Weight, bw.Comment)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the row from the DB
func (bw *Bodyweight) Delete(conn *sql.DB) error {
	// Lookup record by timestamp
	err := conn.QueryRow(`
	DELETE FROM metrics.bodyweight
	WHERE timestamp=$1`,
		bw.Timestamp).Scan(bw)
	if err != nil {
		return err
	}
	return nil
}

/*
	RESTfulHandler
*/

// HandlePost creates a new bodyweight record.
func (bw *Bodyweight) HandlePost(w http.ResponseWriter, req *http.Request) {
	err := torque.ReadBodyTo(w, req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = bw.Create(torque.DBConn); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Created %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// HandleGet returns the related bodyweight record
func (bw *Bodyweight) HandleGet(w http.ResponseWriter, req *http.Request) {
	timestamp, err := torque.GetOrCreateTimestamp(req)
	if err != nil {
		http.Error(w, "Invalid timestamp provided", http.StatusBadRequest)
		return
	}
	bw.Timestamp = timestamp
	if err = bw.Retrieve(torque.DBConn); err != nil {
		http.NotFound(w, req)
		return
	}
	log.Printf("Retrieved %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// HandlePut updates a Bodyweight resource.
func (bw *Bodyweight) HandlePut(w http.ResponseWriter, req *http.Request) {
	// Parse body of PUT request into a Bodyweight struct
	err := torque.ReadBodyTo(w, req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = bw.Update(torque.DBConn); err != nil {
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
	timestamp, err := torque.GetOrCreateTimestamp(req)
	if err != nil {
		http.Error(w, "Invalid timestamp provided", http.StatusBadRequest)
		return
	}
	if err = bw.Delete(torque.DBConn); err != nil {
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
	// user := GetUserFromRequest
	user := users.NewUserAuth()
	return torque.SlashJoin(user.Username, "bodyweight")
}

// RegisterURL sets up the handler for the Bodyweight reosurce on the server.
func (bw *Bodyweight) RegisterURL() error { return nil }
