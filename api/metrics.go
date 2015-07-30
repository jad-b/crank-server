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

// Bodyweight is a timestamped bodyweight record.
type Bodyweight struct {
	RESTfulRouter `json:"-"`
	Bodyweight    float32   `json:"bodyweight"`
	Timestamp     time.Time `json:"timestamp"`
}

// Delete removes the bodyweight record from the database.
func (bw *Bodyweight) Delete(w http.ResponseWriter, req *http.Request) {
	timestamp, err := web.stamp(req)
	// Lookup record by timestamp
	log.Printf("Looking up Bodyweight record from %s from DB", timestamp)
	bw := nil
	if bw == nil {
		http.NotFound(w, req)
		return
	}
	writeOkayJSON(w, bw)
}

// Get returns the related bodyweight record
func (bw *Bodyweight) Get(w http.ResponseWriter, req *http.Request) {
	var bwRec *Bodyweight
	timestamp, err := web.stamp(req)
	log.Printf("Looking up Bodyweight record from %s from DB", timestamp)
	err = db.QueryRow("SELECT bodyweight, timestamp FROM Bodyweight WHERE timestamp = $1",
		timestamp).Scan(bwRec)
	if err != nil {
		log.Printf("Problem reading from database: %s", err.Error())
		http.NotFound(w, req)
	}
	writeOkayJSON(w, bwRec)
}

// Post creates a new bodyweight record.
func (bw *Bodyweight) Post(w http.ResponseWriter, req *http.Request) {
	var bwRec *Bodyweight
	err := web.ReadBodyTo(w, req, bwRec)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	// TODO Create a record in the database
	log.Printf("Creating bodyweight record of %+v", bwRec)
	writeOkayJSON(w, bwRec)
}

// Put updates a Bodyweight resource.
func (bw *Bodyweight) Put(w http.ResponseWriter, req *http.Request) {
	var bwRec *Bodyweight
	ts := web.Stamp(req)
	// TODO Update record from the database using timestamp
	// Write updated record to client
	writeOkayJSON(w, bwRec)
}
