package api

import (
	"log"
	"net/http"
	"time"

	"github.com/jad-b/crank/crank"
)

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

func timeFromQuery(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	if &queryTime == nil {
		log.Print("Failed to retrieve a timestamp from the request")
	}
	return time.Parse(time.RFC3339, queryTime)
}
