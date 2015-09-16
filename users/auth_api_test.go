// +build test db,api

package users

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jad-b/torque"
)

var serverHost = "localhost"

func TestAuthentication(t *testing.T) {
	username, password := "ValidUser", "GreatPassword"
	db := torque.Connect()
	defer db.Close()
	// Create user account for testing
	u := NewUserAccount(username, password)
	if err := u.Create(db); err != nil {
		t.Fatal(err)
	}
	defer u.Delete(db)
	t.Logf("User before authentication:\n%s", torque.PrettyJSON(u))

	// Create a request for authentication
	req, err := buildAuthenticationRequest(serverHost, username, password)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()

	// See how we handle the authentication
	HandleAuthentication(resp, req)

	if resp.HeaderMap.Get(torque.HeaderAuthorization) == "" {
		t.Fatalf("Authorization header not set or empty:\n%s",
			torque.DumpRecordedResponse(resp))
	}
	// White-box: Check database for updated user row
	if err := u.Retrieve(db); err != nil {
		t.Fatal(err)
	}
	nilTime := time.Time{}
	if u.CurrentToken == "" ||
		u.TokenCreated == nilTime ||
		u.TokenLastUsed == nilTime {
		t.Fatal("User row not updated for token authentication:\n%#v", u)
	}
	t.Logf("User after authentication:\n%s", torque.PrettyJSON(u))
}

// Try to authenticate a non-existent User account against the Authentication handler.
func TestBadAccountAuthentication(t *testing.T) {
	db := torque.Connect()
	defer db.Close()
	username, password := "InvalidUser", "AndTheirLamePassword"
	u := &UserAuth{Username: username}
	u.Delete(db)
	// Create a request for authentication
	req, err := buildAuthenticationRequest(serverHost, username, password)
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
	username, password := "EzekielSparks", "Tungsten"
	// Create a request for authentication
	req, err := buildAuthenticationRequest(serverHost, username, password)
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
