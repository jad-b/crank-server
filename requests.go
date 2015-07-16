package http

import (
	"io/ioutil"
	"net/http"
)

// ReadBody extracts the body from the HTTP request. If there is an error, it
// writes it back to the response.
func ReadBody(w http.ResponseWriter, req *http.Request) (b []byte) {
	if b, err := ioutil.ReadAll(req.Body); err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return nil
	}
	return b
}
