// +build test db

package users

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

var (
	username = "EzekielSparks"
	password = "BuildABetterButterRobot"
	user     = NewUserAccount(username, password)
	testURL  = url.URL{
		Scheme: "http",
		Host:   "localhost:18000",
		// Generate a random user ID for working with
		Path: "users",
	}
	db *sqlx.DB
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// Try to authenticate a non-existent User account against the Authentication handler.
func TestBadAccountAuthentication(t *testing.T) {
	serverURL, username, password := "https://localhost", "JohnFritz", "gazebo"
	// Create a request for authentication
	req, err := buildAuthenticationRequest(serverURL, username, password)
	resp := httptest.NewRecorder()

	// See how we handle the authentication
	HandleAuthentication(resp, req)

	// That account doesn't exist - we should see an error in the status code
	// and returned body
	if resp.Code != http.StatusUnauthorized {
		t.Error("Didn't receive StatusUnauthorized for non-existent account")
	}
	var errResp torque.ErrorResponse
	err = torque.ReadJSONResponse(resp.Body, &errResp)
	if err != nil { // Bad response was returned
		t.Errorf("Failed to read 401 response: %s", err.Error())
	}
	if errResp.Error() == "" { // No 'error' message returned
		t.Error("No 'error' value found in response")
	}
}

// Create a new account via the REST API
func TestAccountCreation(t *testing.T) {
	// Create request
	req, err := http.NewRequest("POST", testURL.String(), nil)
	// Set whose account to make by piggy-backing on the Auth
	req.SetBasicAuth(username, password)
	resp := httptest.NewRecorder()

	u := &UserAuth{}
	u.HandlePost(resp, req)

	if resp.Code != 200 {
		t.Error("Account creation did not succeed")
	}
	var respUser UserAuth
	if err = torque.ReadJSONResponse(resp.Body, &respUser); err != nil {
		t.Errorf("Failed to retrieve user from response; %s", err.Error())
	}
}

func TestNewUserAccount(t *testing.T) {
	u := NewUserAccount(username, password)
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
