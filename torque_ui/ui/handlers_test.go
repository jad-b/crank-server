package ui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jad-b/torque"
)

var sampleWktJSON = []byte(`{
  "workout_id": -1,
  "last_modified": "2015-12-28T16:04:08.668Z",
  "user_id": -1,
  "exercises": [
    {
      "exercise_id": 0,
      "workout_id": -1,
	  "last_modified": "2015-12-28T16:04:08.668Z",
      "movement": "Bench Press",
      "modifiers": "",
      "sets": "1 x 1, 2 x 3, 3 x 4",
      "tags": "tag1=value1; tag2=value2"
    },
    {
      "exercise_id": 1,
      "workout_id": -1,
	  "last_modified": "2015-12-28T16:04:08.668Z",
      "movement": "Squat",
      "modifiers": "",
      "sets": "1 x 1, 2 x 3, 3 x 4",
      "tags": "tag1=value1; tag2=value2"
    }
  ],
  "tags": ""
}`)

func TestWorkoutParsing(t *testing.T) {
	var wkt Workout
	if err := json.Unmarshal(sampleWktJSON, &wkt); err != nil {
		t.Fatal(err)
	}
	t.Log(torque.PfmtJSON(&wkt))
}

func TestWktHandlerParsing(t *testing.T) {
	body := bytes.NewBuffer(sampleWktJSON)
	req, err := http.NewRequest("POST", "http://localhost/wkt", body)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	WktHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Got back a %d", w.Code)
	}
}
