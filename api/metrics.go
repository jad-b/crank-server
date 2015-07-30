package api

/*
	URL Path: /metrics
*/

import (
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
	RESTfulRouter `json:"-"`
	Timestamp     time.Time `json:"timestamp"`
	Weight        float32   `json:"weight"`
	Comment       string    `json:"comment"`
}

// Delete removes the bodyweight record from the database.
func (bw *Bodyweight) Delete(w http.ResponseWriter, req *http.Request) {
	// Retrieve timestamp from request
	timestamp, err := web.stamp(req)

	// Lookup record by timestamp
	err := DBConn.QueryRow(`
	DELETE FROM metrics.bodyweight
	WHERE timestamp=$1`,
		timestamp).Scan(bw)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	log.Printf("Deleted bodyweight @ %s", timestamp)
	writeOkayJSON(w, bw)
}

// Get returns the related bodyweight record
func (bw *Bodyweight) Get(w http.ResponseWriter, req *http.Request) {
	timestamp, err := web.stamp(req)
	log.Printf("Looking up Bodyweight record from %s from DB", timestamp)
	err = DBConn.QueryRow(`
	SELECT (timestamp, weight, comment)
	FROM metrics.bodyweight
	WHERE timestamp=$1`,
		timestamp).Scan(bw)
	if err != nil {
		log.Printf("Problem reading from database: %s", err.Error())
		http.NotFound(w, req)
	}
	log.Printf("Retrieved %+v", bw)
	writeOkayJSON(w, bwg)
}

// Post creates a new bodyweight record.
func (bw *Bodyweight) Post(w http.ResponseWriter, req *http.Request) {
	err := web.ReadBodyTo(w, req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	result, err := DBConn.Exec(`
	INSERT INTO metrics.bodyweight (timestamp, weight, comment)
	VALUES ($1, $2, $3)`,
		bw.Timestamp, bw.Weight, bw.Comment)
	if err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
	}
	log.Printf("Created %+v", bw)
	writeOkayJSON(w, bw)
}

// Put updates a Bodyweight resource.
func (bw *Bodyweight) Put(w http.ResponseWriter, req *http.Request) {
	// Parse body of PUT request into a Bodyweight struct
	err := web.ReadBodyTo(w, req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	// Update record in database
	// TODO Only overwrite with provided fields. Maybe by building the SQL
	// statement string w/ conditional logic?
	result, err := DBConn.Exec(`
	UPDATE metrics.bodyweight
	SET weight=$2, comment='$3'
	WHERE timestamp > $1`,
		bw.Timestamp, bw.Weight, bw.Comment)
	i
	if err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
	}
	log.Printf("Updated %+v", bw)
	// Write updated record to client
	writeOkayJSON(w, bw)
}
