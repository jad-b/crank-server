package api

import (
	"encoding/json"
	"net/http"
)


// writeJSON writes the value v to the http response stream as json with standard
// json encoding.
// Stole from github.com/docker/docker/api/server
func writeJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
