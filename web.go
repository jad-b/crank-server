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
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return nil
	}
	return b
}

// ReadBodyTo reads the body of a request directly into a given struct.
func ReadBodyTo(w http.ResponseWriter, req *http.Request, v interface{}) error {
	return json.NewDecoder(req.Body).Decode(v)
}

// stamp ensures a timestamp is attached to the Request. First it looks for
// a Query field "timestamp". Failing that, it returns the current time.
// Query.
func stamp(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	if &queryTime == nil {
		return time.Now(), nil
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
		log.Printf("Failed to encode as json:\n\t%v\nSending %d", v,
			http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeOkayJSON(w http.ResponseWriter, v interface{}) {
	return writeJSON(w, http.StatusOK, v)
}
