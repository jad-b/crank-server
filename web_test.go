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
	type StampStruct struct {
		Timestamp time.Time
	}
	// Test request handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		ts, err := GetTimestampQuery(r)
		if err != nil {
			HTTPError(w, err, http.StatusBadRequest)
			return
		}
		t.Logf("Received: %s", ts.String())
		WriteOkayJSON(w, StampStruct{Timestamp: ts})
	}
	resp := httptest.NewRecorder()
	// setup timestamp in request
	stamp := time.Now().UTC()
	u := &url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   "/timestamp",
	}
	SetTimestampQuery(u, stamp)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Query string: %s", req.URL.RawQuery)

	handler(resp, req)

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
