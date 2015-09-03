package users

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

var (
	username = "EzekielSparks"
	password = "BuildABetterButterRobot"
	testURL  = url.URL{
		Scheme: "http",
		Host:   "localhost:18000",
		// Generate a random user ID for working with
		Path: "users",
	}
	DBConn *sqlx.DB
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

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

func TestCRUDExercise(t *testing.T) {
	// Setup our database connection
	pgConf := torque.LoadPostgresConfig(*torque.PsqlConf)
	DBConn = torque.OpenDBConnection(pgConf)
	defer DBConn.Close()

	// Helper logging function
	cry := func(msg string, e error) {
		t.Fatalf("%s: %s", msg, e.Error())
	}

	// Create the table, if missing
	_, err := DBConn.Exec(UserAuthSQL)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		cry("Failed to create UserAuth table", err)
	}

	// Create a non-existent record
	u := NewUserAccount(username, password)
	// Try to clean-up; Delete may not work
	defer (func() {
		u.Delete(DBConn)
	})()

	// Try and retrieve the record before it exists
	err = u.Retrieve(DBConn)
	if err == nil {
		cry("Test user already exists; aborting", err)
	}

	// Create the record
	err = u.Create(DBConn)
	if err != nil {
		cry("Failed to create record", err)
	}
	// Retrieve newly-created record
	u2 := &UserAuth{Username: username}
	if err = u2.Retrieve(DBConn); err != nil {
		cry("Failed to retrieve record", err)
	}
	// They should look the same - somewhat
	if u.PasswordHash != u2.PasswordHash {
		t.Fatal("Failed to retrieve correct account - password hashes don't match")
	}

	// Try to update the record via stamping a new token
	err = u2.Authorize(DBConn)
	if err != nil {
		cry("Failed to update the record during Authorization", err)
	}
	// Retrieve changes
	err = u.Retrieve(DBConn)
	if err != nil {
		cry("Failed to retrieve newly-updated record", err)
	}
	if u.CurrentToken != u2.CurrentToken {
		cry("Token doesn't match after update", err)
	}

	// Delete the user record
	err = u2.Delete(DBConn)
	if err != nil {
		cry("Failed to delete record", err)
	}
	// Should fail to retrieve a deleted record
	err = u.Retrieve(DBConn)
	if err == nil {
		cry("Post-deletion retrieval returned SUCCESS", err)
	}
}
