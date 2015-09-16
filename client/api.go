package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/users"
)

// TorqueAPI is a client-side representation of a Torque server connection.
type TorqueAPI struct {
	ServerURL url.URL        `json:"server_url"`
	User      users.UserAuth `json:"user"`
}

// NewTorqueAPI instantiates a new API connection from a URL string.
func NewTorqueAPI(serverURL string) *TorqueAPI {
	u, err := url.Parse(serverURL)
	if err != nil {
		// No point in continuing if we can't connect to the server
		log.Fatal(err)
	}
	// TODO(jdb) Set up HTTPS certs
	return &TorqueAPI{ServerURL: *u}
}

// Authenticate logs the User in on the Torque Server.
// This is a client-side call.
func (t *TorqueAPI) Authenticate(username, password string) error {
	req, err := users.buildAuthenticationRequest(t.ServerURL.Host, username, password)
	if err != nil {
		return err
	}
	// Send the auth request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("No response received from authentication request")
	}
	if resp.StatusCode != 200 { // Invalid creds
		torque.LogResponse(resp)
		var errResp torque.ErrorResponse
		err = torque.ReadJSONResponse(resp.Body, &errResp)
		if err != nil {
			return errors.New("Failed to read authentication response body")
		}
		return errResp
	}

	// Parse the response into a User object
	err = torque.ReadJSONResponse(resp.Body, &t.User)
	if err != nil {
		return err
	}
	return nil
}

// PostJSON is a convenience wrapper for common POST functionality. This
// includes setting the content-type to "application/json", and marshalling
// structs into JSON.
func PostJSON(serverURL string, res torque.RESTfulResource) (resp *http.Response, err error) {
	payload, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	// Attach our resource to the URL
	postURL := strings.Join([]string{serverURL, res.GetResourceName()}, "/")
	return http.Post(postURL, "application/json", bytes.NewBuffer(payload))
}

// PrepareGetURL converts the
func PrepareGetURL(serverURL string, res torque.RESTfulResource) (*url.URL, error) {
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
func BuildResourcePath(serverURL string, res torque.RESTfulResource) string {
	return strings.Join([]string{serverURL, res.GetResourceName()}, "/")
}
