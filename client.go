package torque

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// A RESTfulClient performs basic data operations over HTTP to a given URL.
//
// It should probably take a variable set of arguments, or a struct defining
// the set of possible options. But that's next, not now.
type RESTfulClient interface {
	RESTfulResource
	HTTPPost(serverURL string) (*http.Response, error)
	HTTPGet(serverURL string) (*http.Response, error)
	HTTPPut(serverURL string) (*http.Response, error)
	HTTPDelete(serverURL string) (*http.Response, error)
}

// API is a client-side representation of a Torque server connection.
type API struct {
	ServerURL url.URL `json:"server_url"`
}

// NewTorqueAPI instantiates a new API connection from a URL string.
func NewTorqueAPI(serverURL string) {
	u, err := url.Parse(serverURL)
	if err != nil {
		// No point in continuing if we can't connect to the server
		log.Fatal(err)
	}
	// TODO(jdb) Set up HTTPS certs
	return API{ServerURL: u}
}

// PostJSON is a convenience wrapper for common POST functionality. This
// includes setting the content-type to "application/json", and marshalling
// structs into JSON.
func PostJSON(serverURL string, res RESTfulResource) (resp *http.Response, err error) {
	payload, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	// Attach our resource to the URL
	postURL := strings.Join([]string{serverURL, res.GetResourceName()}, "/")
	return http.Post(postURL, "application/json", bytes.NewBuffer(payload))
}

// PrepareGetURL converts the
func PrepareGetURL(serverURL string, res RESTfulResource) (*url.URL, error) {
	// Turn the base URL into a safer working form; url.URL
	u, err := url.Parse(serverURL)
	if err != nil {
		return &url.URL{}, err
	}
	// Add our resource's endpoint
	u.Path = strings.Join([]string{serverURL, res.GetResourceName()}, "/")
	return u, nil
}

// NewJSONRequest builds an http.Request with a content-type of
// application/json
func NewJSONRequest(method, serverURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, serverURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// BuildResourcePath builds the absolute URL for the UserAuth resource
func BuildResourcePath(serverURL string, res RESTfulResource) string {
	return strings.Join([]string{serverURL, res.GetResourceName()}, "/")
}
