package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/users"
)

// HTTPGetter tells *you* how it's gonna be GETted
type HTTPGetter interface {
	Get(url.URL) url.URL
}

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
	req, err := t.NewRequest("POST", postURL, res)
	if err != nil {
		return nil, err
	}
	return t.Client.Do(req)
}

// Get retrieves a resource from the Torque server.
func (t *TorqueAPI) Get(res torque.RESTfulResource, params url.Values) (resp *http.Response, err error) {
	var getURL url.URL
	if getter, ok := interface{}(res).(HTTPGetter); ok {
		log.Print("Delegating GET to resource")
		// Delegate URL customization to resource
		SetUserID(res, t.User.ID)
		getURL = getter.Get(*t.BuildURL(res, nil))
	} else { // Build the GET URL ourselves
		getURL = *t.BuildURL(res, params)
	}
	log.Printf("GET URL: %s", getURL.String())
	// Create Request w/ req'd headers
	req, err := t.NewRequest("GET", getURL.String(), nil)
	if err != nil {
		return nil, err
	}
	log.Print("Built GET request")
	torque.LogRequest(req)
	return t.Client.Do(req)
}

// Put updates a resource on the Torque server.
func (t *TorqueAPI) Put(res torque.RESTfulResource) (resp *http.Response, err error) {
	putURL := t.BuildURL(res, nil).String()
	req, err := t.NewRequest("PUT", putURL, res)
	if err != nil {
		return nil, err
	}
	return t.Client.Do(req)
}

// Delete retrieves a resource from the Torque server.
// You may provide JSON to pass options to the server.
func (t *TorqueAPI) Delete(res torque.RESTfulResource) (resp *http.Response, err error) {
	deleteURL := t.BuildURL(res, nil).String()
	req, err := t.NewRequest("DELETE", deleteURL, res)
	if err != nil {
		return nil, err
	}
	return t.Client.Do(req)
}

// NewRequest prepares a new HTTP request.
// It handles filling in the appropriate auth fields.
func (t *TorqueAPI) NewRequest(method string, url string, body interface{}) (*http.Request, error) {
	var payload []byte
	var err error
	if body != nil { // Try and set missing UserID fields
		SetUserID(body, t.User.ID)
		payload, err = json.Marshal(body) // Marshal body into JSON
		if err != nil {
			return nil, err
		}
		log.Print(torque.PrettyJSON(body))
	}
	// Create HTTP Request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	// Set headers
	t.authRequest(req)
	torque.LogRequest(req)
	return req, nil
}

func (t *TorqueAPI) authRequest(req *http.Request) {
	req.Header.Set(torque.HeaderContentType, torque.MimeJSON)
	if t.User.IsAuthenticated() { // Only set if valid
		req.Header.Set(users.AuthHeader(&t.User))
	} else {
		// TODO Perform re-authentication on behalf of user
		log.Print("User's authentication is either missing or invalid")
	}
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

// SetUserID attaches the active user's ID to the resource.
// If a 'UserID' field is found set, it won't do a thing.
func SetUserID(v interface{}, userID int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Failed to set UserID: %s\n", r)
		}
	}()

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	uIDField := val.FieldByName("UserID")
	zeroValue := reflect.Value{}
	if uIDField == zeroValue { // No UserID field was found
		log.Printf("No UserID field in %s", val.Type().String())
	} else if uIDField.Kind() != reflect.Int { // It's not an Integer
		log.Print("UserID is %s, not Integer", uIDField.Type().String())
	} else if uIDField.Int() != 0 { // UserID's been set.
		log.Printf("UserID field has already been set to %d", uIDField.Int())
	} else {
		uIDField.SetInt(int64(userID))
	}
	return nil
}
