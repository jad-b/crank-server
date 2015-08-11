package metrics

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/torque"
)

// Bodyweight is a timestamped bodyweight record, with optional comment.
//
// SQL:
// CREATE TABLE bodyweight (
//	 "timestamp" timestamp(0) with time zone NOT NULL UNIQUE,
//	 weight numeric(5,2) NOT NULL CHECK (weight < 1000),
//	 comment text
// );
type Bodyweight struct {
	Timestamp time.Time `json:"timestamp"`
	Weight    float64   `json:"weight"`
	Comment   string    `json:"comment"`
}

/*
	DBResourcer
*/

// DBCreate inserts a new bodyweight entry into the DB.
func (bw *Bodyweight) DBCreate(conn *sql.DB) error {
	_, err := conn.Exec(`
	INSERT INTO metrics.bodyweight (timestamp, weight, comment)
	VALUES ($1, $2, $3)`,
		bw.Timestamp, bw.Weight, bw.Comment)
	if err != nil {
		return err
	}
	return nil
}

// DBRetrieve does a lookup for the corresponding bodyweight record by timestamp.
func (bw *Bodyweight) DBRetrieve(conn *sql.DB) error {
	log.Printf("Looking up Bodyweight record from %s from DB", bw.Timestamp)
	err := conn.QueryRow(`
	SELECT (timestamp, weight, comment)
	FROM metrics.bodyweight
	WHERE timestamp=$1`,
		bw.Timestamp).Scan(bw)
	if err != nil {
		log.Printf("Problem reading from database: %s", err.Error())
		return err
	}
	return nil
}

// DBUpdate modifies the matching row in the DB by timestamp.
func (bw *Bodyweight) DBUpdate(conn *sql.DB) error {
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

// DBDelete removes the row from the DB
func (bw *Bodyweight) DBDelete(conn *sql.DB) error {
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

// Post creates a new bodyweight record.
func (bw *Bodyweight) Post(w http.ResponseWriter, req *http.Request) {
	err := torque.ReadBodyTo(w, req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = bw.DBCreate(torque.PGConn); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Created %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// Get returns the related bodyweight record
func (bw *Bodyweight) Get(w http.ResponseWriter, req *http.Request) {
	timestamp, err := torque.Stamp(req)
	if err != nil {
		http.Error(w, "Invalid timestamp provided", http.StatusBadRequest)
		return
	}
	bw.Timestamp = timestamp
	if err = bw.DBRetrieve(torque.PGConn); err != nil {
		http.NotFound(w, req)
		return
	}
	log.Printf("Retrieved %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// Put updates a Bodyweight resource.
func (bw *Bodyweight) Put(w http.ResponseWriter, req *http.Request) {
	// Parse body of PUT request into a Bodyweight struct
	err := torque.ReadBodyTo(w, req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = bw.DBUpdate(torque.PGConn); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Updated %+v", bw)
	// Write updated record to client
	torque.WriteOkayJSON(w, bw)
}

// Delete removes the bodyweight record from the database.
func (bw *Bodyweight) Delete(w http.ResponseWriter, req *http.Request) {
	// Retrieve timestamp from request
	timestamp, err := torque.Stamp(req)
	if err != nil {
		http.Error(w, "Invalid timestamp provided", http.StatusBadRequest)
		return
	}
	if err = bw.DBDelete(torque.PGConn); err != nil {
		http.NotFound(w, req)
		return
	}
	log.Printf("Deleted bodyweight @ %s", timestamp)
	torque.WriteOkayJSON(w, bw)
}

/*
	FlagParser
*/

// ParseFlags handles command-line argument parsing.
func (bw *Bodyweight) ParseFlags(action string, args []string) error {
	// Define sub-flags for the bodyweight resource
	var tsFlag torque.TimestampFlag
	bwFlags := flag.NewFlagSet("bwFlags", flag.ContinueOnError)
	bwFlags.Var(&tsFlag, "timestamp", "")
	weight := bwFlags.Float64("weight", 0.0, "")
	comment := bwFlags.String("comment", "", "")

	// Parse the given flags
	bwFlags.Parse(args)
	bw = &Bodyweight{time.Time(tsFlag), *weight, *comment}

	switch action {
	case "create":
		return bw.DBCreate(torque.PGConn)
	case "retrieve":
		return bw.DBRetrieve(torque.PGConn)
	case "update":
		return bw.DBUpdate(torque.PGConn)
	case "delete":
		return bw.DBDelete(torque.PGConn)
	default:
		log.Fatalf("%s is an invalid action", action)
		return nil
	}
}

/*
	RESTfulClient
*/

// ClientPOST creates a new bodyweight record on the REST API server.
//
// It probably makes more sense to have a generic 'metrics/' endpoint that accepts
// a variety of metrics, especially if these continue to grow.
func (bw *Bodyweight) ClientPOST() (resp *http.Response, err error) {
	endpoint := "/metrics/bodyweight" // For now.
	return torque.PostJSON(endpoint, bw)
}
