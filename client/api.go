package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/users"
)

// HTTPPoster tells *you* how it's gonna POST
type HTTPPoster interface {
	Post(url.URL) (*http.Response, error)
}

// TorqueAPI is a client-side representation of a Torque server connection.
type TorqueAPI struct {
	http.Client `json:"-"`
	ServerURL   url.URL        `json:"server_url"`
	User        users.UserAuth `json:"user"`
}

// NewTorqueAPI instantiates a new API connection from a URL string.
func NewTorqueAPI(serverURL string) *TorqueAPI {
	u, err := url.Parse(serverURL)
	// Overwrite whatever was given with what Torque enforces (https)
	u.Scheme = torque.Scheme
	if err != nil {
		// No point in continuing if we can't connect to the server
		log.Fatal(err)
	}
	return &TorqueAPI{ServerURL: *u}
}

// Authenticate logs the User in on the Torque Server.
// This is a client-side call.
func (t *TorqueAPI) Authenticate(username, password string) error {
	req, err := users.BuildAuthenticationRequest(t.ServerURL.Host, username, password)
	if err != nil {
		return err
	}
	// Send the auth request
	resp, err := t.Do(req)
	if err != nil {
		return errors.New("No response received from authentication request")
	}
	if resp.StatusCode != 200 { // Invalid creds
		torque.LogResponse(resp)
		var errResp torque.ErrorResponse
		err = torque.ReadJSONResponse(resp, &errResp)
		if err != nil {
			return errors.New("Failed to read authentication response body")
		}
		return errResp
	}

	// Parse the response into a User object
	err = torque.ReadJSONResponse(resp, &t.User)
	if err != nil {
		return err
	}
	log.Printf("Authenticated to Torque API as %s(%d)\n", username, t.User.ID)
	return nil
}

// Post is a convenience wrapper for common POST functionality.
func (t *TorqueAPI) Post(res torque.RESTfulResource) (resp *http.Response, err error) {
	postURL := t.BuildURL(res, nil).String()
	req, err := t.NewRequest("POST", postURL, res, nil)
	if err != nil {
		return nil, err
	}
	return t.Client.Do(req)
}

// Get retrieves a resource from the Torque server.
func (t *TorqueAPI) Get(res torque.RESTfulResource, params url.Values) (resp *http.Response, err error) {
	getURL := t.BuildURL(res, params).String()
	return t.Client.Get(getURL)
}

// Put updates a resource on the Torque server.
func (t *TorqueAPI) Put(res torque.RESTfulResource) (resp *http.Response, err error) {
	putURL := t.BuildURL(res, nil).String()
	req, err := t.NewRequest("PUT", putURL, res, nil)
	if err != nil {
		return nil, err
	}
	return t.Client.Do(req)
}

// Delete retrieves a resource from the Torque server.
// You may provide JSON to pass options to the server.
func (t *TorqueAPI) Delete(res torque.RESTfulResource, body interface{}) (resp *http.Response, err error) {
	deleteURL := t.BuildURL(res, nil).String()
	req, err := t.NewRequest("DELETE", deleteURL, body, nil)
	if err != nil {
		return nil, err
	}
	return t.Client.Do(req)
}

// NewRequest prepares a new HTTP request.
// It handles filling in the appropriate auth fields.
func (t *TorqueAPI) NewRequest(method string, url string, body interface{}, params url.Values) (*http.Request, error) {
	var payload []byte
	var err error
	// Marshal body into JSON
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	// Create HTTP Request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	// Set headers
	req.Header.Set(torque.HeaderContentType, torque.MimeJSON)
	if t.User.IsAuthenticated() { // Only set if valid
		req.Header.Set(users.AuthHeader(&t.User))
	} else {
		return nil, fmt.Errorf("%s is not authenticated", t.User.Username)
	}
	return req, nil
}

// BuildURL creates a full-fledged URL
func (t *TorqueAPI) BuildURL(res torque.RESTfulResource, params url.Values) *url.URL {
	// Copy the URL
	earl := t.ServerURL
	// Set query parameters, if they exist
	if params != nil {
		earl.RawQuery = params.Encode()
	}
	// Build the API resource path
	earl.Path = t.BuildPath(res)
	return &earl
}

// BuildPath builds the resource path
func (t *TorqueAPI) BuildPath(res torque.RESTfulResource) string {
	return torque.SlashJoin(t.ServerURL.Path, res.GetResourceName())
}

// String pretty-prints the Torque API client.
func (t *TorqueAPI) String() string {
	return torque.PrettyJSON(t)
}
