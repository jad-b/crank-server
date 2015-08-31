package users

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jad-b/torque"
)

var (
	username = "EzekielSparks"
	password = "BuildABetterButterRobot"
)

func init() {
	// Setup our database connection
	pgConf := torque.LoadPostgresConfig()
	DBConn := torque.OpenDBConnection(pgConf)

	testURL := url.URL{
		Scheme: "http",
		Host:   "localhost:18000",
		// Generate a random user ID for working with
		Path: "users",
	}
}

func TestBadAccountAuthentication(t *testing.T) {
	serverURL, username, password := "https://localhost", "JohnFritz", "gazebo"
	// Create a request for authentication
	req, err := buildAuthenticationRequest(serverURL, username, password)
	w := httptest.NewRecorder()

	// See how we handle the authentication
	HandleAuthentication(w, req)

	// That account doesn't exist - we should see an error in the status code
	// and returned body
	if w.Code != http.StatusUnauthorized {
		t.Error("Didn't receive StatusUnauthorized for non-existent account")
	}
	var errResp torque.ErrorResponse
	err := torque.ReadJSONResponse(resp, &errResp)
	if err != nil { // Bad response was returned
		t.Errorf("Failed to read 401 response: %s", err.Error())
	}
	if errResp.Error == "" { // No 'error' message returned
		t.Error("No 'error' value found in response")
	}
}

func TestAccountCreation(t *testing.T) {
	// Create request
	req, err := http.NewRequest("POST", testURL.String(), nil)
	// Set whose account to make by piggy-backing on the Auth
	req.SetBasicAuth(username, password)
	resp := httptest.NewRecorder()

	u := *UserAuth{}
	u.HandlePost(resp, req)

	if resp.Code != 200 {
		t.Error("Account creation did not succeed")
	}
	if u.PasswordHash == "" {
		t.Error("No password was created")
	}
	if u.PasswordHash == password {
		t.Error("Password wasn't hashed")
	}
	if u.PasswordSalt == "" {
		t.Error("No salt created for account")
	}
	if u.Cost == 0 {
		t.Error("No hashing-cost set on account")
	}
}
