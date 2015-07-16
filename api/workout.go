package api

/*
	/workout
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetWorkoutHandler returns a workout by timestamp
func GetWorkoutHandler(w http.ResponseWriter, req *http.Request) {
	timestamp, err := Torque.timeFromQuery(req)
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
