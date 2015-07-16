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
	Bodyweight float32 `json:"bodyweight"`
	// 'omitempty' => Skip the timestamp if it's empty
	Timestamp time.Time `json:"timestamp"`
}

// DeleteBodyweight removes the bodyweight record from the database.
func DeleteBodyweight(w http.ResponseWriter, req *http.Request) {
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

// GetBodyweight returns the related bodyweight record
func GetBodyweight(w http.ResponseWriter, req *http.Request) {
	timestamp, err := web.stamp(req)
	// Lookup bodyweight from DB
	log.Printf("Looking up Bodyweight record from %s from DB", timestamp)
	bw := nil
	if bw == nil {
		http.NotFound(w, req)
		return
	}
	writeOkayJSON(w, bw)
}

// PostBodyweight creates a new bodyweight record.
func PostBodyweight(w http.ResponseWriter, req *http.Request) {
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
