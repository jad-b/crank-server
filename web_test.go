package torque

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestTimestampInRequests(t *testing.T) {
	// Test setup
	type StampStruct struct {
		Timestamp time.Time
	}
	// 3) Out-of-order; see below:
	handler := func(w http.ResponseWriter, r *http.Request) {
		//   a) Server parses timestamp from most-to-least precise
		ts, err := GetTimestampQuery(r)
		if err != nil {
			HTTPError(w, err, http.StatusBadRequest)
			return
		}
		// b) DB op(s) go here...
		t.Logf("Received: %s", ts.String())
		WriteOkayJSON(w, StampStruct{Timestamp: ts})
	}
	resp := httptest.NewRecorder()

	// 1) Client prepares request
	stamp := time.Now().UTC()
	u := &url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   "/timestamp",
	}
	// 1a) Client sets timestamp to appropriate precision
	SetTimestampQuery(u, stamp)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Query string: %s", req.URL.RawQuery)

	// 2) Request is sent to server
	handler(resp, req)

	// Assertions
	if resp.Code != 200 {
		t.Fatal(resp.Body.String())
	}
	// Check body for timestamp
	tsRecvd := StampStruct{}
	err = json.NewDecoder(resp.Body).Decode(&tsRecvd)
	if err != nil {
		t.Fatal(err)
	}
	if Stamp(tsRecvd.Timestamp) != Stamp(stamp) {
		t.Fatalf("%s != %s", Stamp(tsRecvd.Timestamp), Stamp(stamp))
	}
}
