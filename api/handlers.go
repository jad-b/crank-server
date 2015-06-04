package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jad-b/crank/crank"
	"log"
	"net/http"
	"time"
)

// GetWorkoutHandler returns a workout by timestamp
func GetWorkoutHandler(w http.ResponseWriter, req *http.Request) {
	// log.Print("Requested /workout/")
	// Write a stub workout to the reesponse
	// Return the workout of the requested timestamp, if existing.
	buf := new(bytes.Buffer)
	req.Write(buf)
	log.Printf("Server: Received this request:\n%s", buf)

	queryTime := req.URL.Query().Get("timestamp")
	if &queryTime == nil {
		log.Fatal("Failed to retrieve a timestamp from the request")
	}
	log.Printf("Server: request timestamp is %+v", queryTime)
	timestamp, err := time.Parse(time.RFC3339, queryTime)
	if err != nil {
		log.Fatalf("This '%s' couldn't be parsed:\n%s", queryTime, err)
	}
	log.Printf("Server: Workout requested w/ timestamp: %s", timestamp)
	workout, err := crank.LookupWorkout(timestamp)
	if err != nil {
		now := time.Now()
		workout = &crank.Workout{
			Timestamp: now,
			Comment:   fmt.Sprintf("Time is %s", now.String()),
		}
	}
	// log.Printf("Workout: %+v", workout)
	err = json.NewEncoder(w).Encode(workout)
	if err != nil {
		log.Fatal("Failed to write Workout to response")
	}
	log.Print("Server: Sending workout")
	w.Header().Set("Content-Type", "application/json")
}
