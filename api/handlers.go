package api

import (
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
	now := time.Now()
	workout := crank.Workout{
		Timestamp: now,
		Comment:   fmt.Sprintf("Time is %s", now.String()),
	}
	// log.Printf("Workout: %+v", workout)
	err := json.NewEncoder(w).Encode(workout)
	if err != nil {
		log.Fatal("Failed to write Workout to response")
	}
	w.Header().Set("Content-Type", "application/json")
}
