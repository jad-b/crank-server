package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// Resuable Client
var c = &http.Client{}

// PostJSON is a convenience wrapper for common POST functionality. This
// includes setting the content-type to "application/json", and marshalling
// structs into JSON.
func PostJSON(endpoint string, body interface{}) (resp *http.Response, err error) {
	u := url.URL{
		Scheme: "http",
		Host:   "127.0.0.1",
	}
	u.Path = endpoint
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.Post(u.String(), "application/json", bytes.NewBuffer(payload))
}

// CreateBodyweight POSTs a new bodyweight record.
//
// It probably makes more sense to have a generic 'metrics/' endpoint that accepts
// a variety of metrics, especially if these continue to grow.
func CreateBodyweight(bodyweight float32, timestamp time.Time) (resp *http.Response, err error) {
	endpoint := "/metrics/bodyweight" // For now.

	// Not sure if you can even pass 'nil' for time.Time. That would mean we need
	// to accept an 'interface{}' and type assert between nil and a valid
	// time.Time
	if timestamp != nil {
		timestamp = time.Now()
	}

	body := &struct {
		Bodyweight float32 `json:"bodyweight"`
		// 'omitempty' => Skip the timestamp if it's empty
		Timestamp time.Time `json:"timestamp"`
	}{
		Bodyweight: bodyweight,
		Timestamp:  timestamp,
	}
	return PostJSON(endpoint, &body)
}
