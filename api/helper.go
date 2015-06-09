package api

import (
	"encoding/json"
	"log"
	"net/http"
)

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
