package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/crank/crank"
)

// BodyweightRecord encapsulates a record about bodyweight.
type BodyweightRecord struct {
	Bodyweight float32 `json:"bodyweight"`
	// 'omitempty' => Skip the timestamp if it's empty
	Timestamp time.Time `json:"timestamp"`
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
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
	}
	bwRec := &BodyweightRecord{}
	err = json.Unmarshal(body, bwRec)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
	}
	// TODO Create a record in the database
	fmt.Println("Received a bodyweight of %f w/ timestamp %s",
		bwRec.Bodyweight, bwRec.Timestamp)
}

func timeFromQuery(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	if &queryTime == nil {
		log.Print("Failed to retrieve a timestamp from the request")
	}
	return time.Parse(time.RFC3339, queryTime)
}
