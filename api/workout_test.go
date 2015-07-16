package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jad-b/crank/client"
	"github.com/jad-b/crank/crank"
)

const (
	baseURL string = "http://localhost:8000"
)

func TestGetWorkoutHandler(t *testing.T) {
	w := httptest.NewRecorder()
	workoutURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	nowTime := time.Now()
	q := workoutURL.Query()
	q.Set("timestamp", nowTime.Format(time.RFC3339))
	workoutURL.RawQuery = q.Encode()
	log.Printf("Test: Sending %s", workoutURL)
	req, err := http.NewRequest("GET", workoutURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	GetWorkoutHandler(w, req)
	var wkt crank.Workout
	err = json.Unmarshal(w.Body.Bytes(), &wkt)
	if err != nil {
		t.Error(err)
	}
	log.Printf("Test: Returned workout: %+v", wkt)

	if &wkt.Timestamp == nil {
		t.Errorf("Failed to unmarshal workout: %+v", wkt)
	} else if diff := time.Now().Sub(wkt.Timestamp).Seconds(); diff > 1 {
		t.Errorf("Took longer than 100 ms to return workout(%s): %s",
			wkt.Timestamp, diff)
	} else if wkt.Comment != fmt.Sprintf("Time is %s", wkt.Timestamp.Format(time.RFC3339)) {
		t.Error("Fuck, this unit test is really getting specific")
	}
}

func TestGetWorkout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(GetWorkoutHandler))
	defer server.Close()

	testTime := time.Now()
	w, err := client.GetWorkout(server.URL, testTime)
	log.Printf("Test: Received this workout:\n\t%+v", w)
	if err != nil {
		t.Error(err)
	} else if w.Timestamp.Format(time.RFC3339) != testTime.Format(time.RFC3339) {
		t.Errorf("%s != %s", w.Timestamp, testTime)
	} else if w.Comment != fmt.Sprintf("Time is %s", testTime.Format(time.RFC3339)) {
		t.Errorf("Comment was not what we expected")
	}
}
