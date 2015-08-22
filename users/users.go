package users

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/torque"
)

// User is your typical user record. Contains id, personal information,
// password-y fields for authentication, and last logged in time
//
// SQL:
// CREATE TABLE users (
//	 "id" numeric(5,2) NOT NULL
//	 "username" text NOT NULL UNIQUE
//	 "first_name" text NOT NULL
//	 "last_name" text NOT NULL
//	 "email" text NOT NULL UNIQUE
//   "enabled" bool
//	 "password_hash" text NOT NULL
//	 "password_salt" text NOT NULL
//	 "timestamp" timestamp(0) with time zone NOT NULL UNIQUE,
// );
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Enabled   bool      `json:"enabled"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUser creates a new User instance from the provided username
func NewUser(username string) *User {
	return &User{Username: username}
}

/*
	DBActor
*/

// Create inserts a new User into the database
func (u *User) Create(conn *sql.DB) error {
	_, err := conn.Exec(`
	INSERT INTO torque.users (id, username, first_name, last_name, email, enabled, password_hash, password_salt, timestamp)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		u.ID, u.Username, u.FirstName, u.LastName, u.Email, u.Enabled, u.PasswordHash, u.PasswordSalt, time.Now())

	if err != nil {
		return err
	}
	return nil
}

// Retrieve a User from the DB by filtering on Username
func (u *User) Retrieve(conn *sql.DB) error {
	log.Printf("Looking up User from %s from DB", user.Username)
	err := conn.QueryRow(`
	SELECT (id, username, first_name, last_name, email, enabled, password_hash, password_salt, timestamp)
	FROM torque.users
	WHERE username=$1`,
		u.Username).Scan(u)
	if err != nil {
		log.Printf("Problem reading from database: %s", err.Error())
		return err
	}
	return nil
}

// Update a user entry in the database
func (u *User) Update(conn *sql.DB) error {
	_, err := conn.Exec(`
	UPDATE torque.users
	SET first_name='$2', last_name='$3', email='$4', enabled=$5, timestamp=$6
	WHERE username=$1`,
		u.Username, u.FirstName, u.LastName, u.Email, u.Enabled, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a User record from the database. In most cases it will
// probably be best practice to simply flag a user as disabled via a PUT, but
// we do also need to expose this ability.
func (u *User) Delete(conn *sql.DB) error {
	err := conn.QueryRow(`
	DELETE FROM torque.users
	WHERE username=$1`, u.Username).Scan(u)
	if err != nil {
		return nil, err
	}
}

/*
	RESTfulHandler
*/

// HandlePost creates a new User record.
func (u *User) HandlePost(w http.ResponseWriter, req *http.Request) {
	err := torque.ReadBodyTo(w, req, u)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = u.Create(torque.DBConn); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Created %+v", u)
	torque.WriteOkayJSON(w, u)
}

// HandleGet returns the specified User record
func (u *User) HandleGet(w http.ResponseWriter, req *http.Request) {
	timestamp, err := torque.Stamp(req)
	if err != nil {
		http.Error(w, "Invalid timestamp provided", http.StatusBadRequest)
		return
	}
	u.Timestamp = timestamp
	if err = u.Retrieve(torque.DBConn); err != nil {
		http.NotFound(w, req)
		return
	}
	log.Printf("Retrieved %+v", u)
	torque.WriteOkayJSON(w, u)
}

// HandlePut updates a User resource.
func (u *User) HandlePut(w http.ResponseWriter, req *http.Request) {
	// Parse body of PUT request into a User struct
	err := torque.ReadBodyTo(w, req, u)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = u.Update(torque.DBConn); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Updated %+v", u)
	// Write updated record to client
	torque.WriteOkayJSON(w, u)
}

// HandleDelete removes a user record from the database.
func (u *User) HandleDelete(w http.ResponseWriter, req *http.Request) {
	// Retrieve timestamp from request
	timestamp, err := torque.Stamp(req)
	if err != nil {
		http.Error(w, "Invalid timestamp provided", http.StatusBadRequest)
		return
	}
	if err = u.Delete(torque.DBConn); err != nil {
		http.NotFound(w, req)
		return
	}
	log.Printf("Deleted user %s @ %s", u.Username, timestamp)
	torque.WriteOkayJSON(w, u)
}

/*
	RESTfulClient
*/

// HTTPPost creates a new user record on the REST API server.
func (u *User) HTTPPost(serverURL string) (resp *http.Response, err error) {
	endpoint := "/users"
	return torque.PostJSON(endpoint, u)
}

// HTTPGet requests a user record from the server.
func (u *User) HTTPGet(serverURL string) (resp *http.Response, err error) {
	return nil, nil
}

// HTTPPut updates the server with the current state of the User record.
func (u *User) HTTPPut(serverURL string) (resp *http.Response, err error) {
	return nil, nil
}

// HTTPDelete deletes the matching user record on the server.
func (u *User) HTTPDelete(serverURL string) (resp *http.Response, err error) {
	return nil, nil
}
