package api

import (
	"encoding/json"
	"fmt"
	"github.com/jad-b/crank/client"
	"github.com/jad-b/crank/crank"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost:8000"
)

func TestGetWorkoutHandler(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	GetWorkoutHandler(w, req)
	var wkt crank.Workout
	json.Unmarshal(w.Body.Bytes(), &wkt)
	// log.Printf("Returned workout: %+v", wkt)

	if &wkt.Timestamp == nil {
		t.Errorf("Failed to unmarshal workout: %+v", wkt)
	} else if diff := time.Now().Sub(wkt.Timestamp).Seconds(); diff > .1 {
		t.Errorf("Took longer than 100 ms to return workout(%s): %s",
			wkt.Timestamp, diff)
	} else if wkt.Comment != fmt.Sprintf("Time is %s", wkt.Timestamp.String()) {
		t.Error("Fuck, this unit test is really getting specific")
	}
}

func TestGetWorkout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(GetWorkoutHandler))
	defer server.Close()

	testTime := time.Now()
	w, err := client.GetWorkout(server.URL, testTime)
	if err != nil {
		t.Errorf("Error wasn't nil! %s", err)
	} else if w.Timestamp != testTime {
		t.Errorf("Expected the same time to be returned")
	} else if w.Comment != fmt.Sprintf("Time is %s", testTime.String()) {
		t.Errorf("Comment was not what we expected")
	}
}
