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
	log.Print("Requested /workout/")
	// Write a stub workout to the reesponse
	now := time.Now()
	json.NewEncoder(w).Encode(&crank.Workout{
		Timestamp: now,
		Comment:   fmt.Sprintf("Time is %s", now.String()),
	})
}
