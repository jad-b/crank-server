package api

import (
	"fmt"
	"github.com/jad-b/crank/client"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
