package torque

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// A RESTfulClient performs basic data operations over HTTP to a given URL.
//
// It should probably take a variable set of arguments, or a struct defining
// the set of possible options. But that's next, not now.
type RESTfulClient interface {
	Create(serverURL string) *http.Response
	Retrieve(serverURL string) *http.Response
	Update(serverURL string) *http.Response
	Delete(serverURL string) *http.Response
}

// PostJSON is a convenience wrapper for common POST functionality. This
// includes setting the content-type to "application/json", and marshalling
// structs into JSON.
func PostJSON(serverURL string, body interface{}) (resp *http.Response, err error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return http.Post(serverURL, "application/json", bytes.NewBuffer(payload))
}
