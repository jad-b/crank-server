package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/crank/crank"
	"github.com/jad-b/crank/http"
)

// timefromQuery extracts and parses a RFC3339 timestamp from the request
// Query.
func timeFromQuery(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	if &queryTime == nil {
		log.Print("Failed to retrieve a timestamp from the request")
	}
	return time.Parse(time.RFC3339, queryTime)
}

// GetWorkoutHandler returns a workout by timestamp
func GetWorkoutHandler(w http.ResponseWriter, req *http.Request) {
	timestamp, err := timeFromQuery(req)
	workout, err := crank.LookupWorkout(timestamp)
	if err != nil {
		http.NotFound(w, req) // Write 404 to response
		return
	}
	writeJSON(w, http.StatusOK, workout)
}

// PostBodyweightHandler creates a new bodyweight record.
func PostBodyweightHandler(w http.ResponseWriter, req *http.Request) {
	var bwRec *BodyweightRecord
	body := http.ReadBody(w, req)
	if err = json.Unmarshal(body, bwRec); err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
	}
	// TODO Create a record in the database
	fmt.Println("Received a bodyweight of %f w/ timestamp %s",
		bwRec.Bodyweight, bwRec.Timestamp)
}
