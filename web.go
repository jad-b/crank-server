package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// LogHTTPError captures and writes a 500 HTTP error.
func LogHTTPError(w http.ResponseWriter, err error) {
	log.Print(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// ReadBody extracts the body from the HTTP request. If there is an error, it
// writes it back to the response.
func ReadBody(w http.ResponseWriter, req *http.Request) (b []byte) {
	if b, err := ioutil.ReadAll(req.Body); err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return nil
	}
	return b
}

// timefromQuery extracts and parses a RFC3339 timestamp from the request
// Query.
func timeFromQuery(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	if &queryTime == nil {
		log.Print("Failed to retrieve a timestamp from the request")
	}
	return time.Parse(time.RFC3339, queryTime)
}

// writeJSON writes the value v to the http response stream as json with standard
// json encoding.
// Stole from github.com/docker/docker/api/server
func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Print("Failed to encode as json:\n\t%v\nSending %d", v,
			http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
