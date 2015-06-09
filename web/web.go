/*
Package web provides a bunch of helper function for web servicing.
*/
package web

import (
	"log"
	"net/http"
)

// LogHTTPError captures and writes a 500 HTTP error.
func LogHTTPError(w http.ResponseWriter, err error) {
	log.Print(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
