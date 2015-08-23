package users

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/torque"
)

// UserAuthSQL is SQL for creating the auth.UserAuths table
const UserAuthSQL = `
CREATE TABLE auth.UserAuths(
  "id" numeric(5,2) NOT NULL
  "user_id" text NOT NULL UNIQUE
  "password_hash" text NOT NULL
  "password_salt" text NOT NULL
  "iteration_count" numeric(5,2) NOT NULL
  "current_token" text NOT NULL
  "timestamp" timestamp(0) with time zone NOT NULL UNIQUE,
  "token_last_seen" timestamp(0) with time zone NOT NULL UNIQUE,
);`

// UserAuth holds the auth token data mapped to a UserAuth ID.
// TODO It also holds some data, like Enabled, AccountCreation, that might be better
// served in a UserMeta table, but I left here for simplicity's sake. Also, I
// hear JOINs are expensive.
type UserAuth struct {
	// Primary key
	ID int `json:"id"`
	// Because it doesn't live in UserPII
	Username        string    `json:"username"`
	AccountCreation time.Time `json:"account_creation"`
	Enabled         bool      `json:"enabled"`
	// Salt used to hash the password
	PasswordSalt string
	// Type of hash used
	PasswordHash string
	// Power-of-two times we iterated over the stored password when hashing
	Cost int
	// Currently active auth token
	CurrentToken string `json:"token"`
	// What time does this apply to? Token creation?
	Timestamp time.Time `json:"timestamp"`
	// Last time the token was used in an API request
	TokenLastUsed time.Time `json:"token_last_used"`
}

// NewUserAuth creates a new UserAuth instance with some defaults in place.
func NewUserAuth() *UserAuth {
	return &UserAuth{
		PasswordSalt:   NewSalt(DefaultSaltLength),
		IterationCount: DefaultIterationCount,
		Timestamp:      time.Now(),
	}
}

// ValidatePassword verifies if the given username/password is valid.
func (u *UserAuth) ValidatePassword(password string) bool {
	// Hash the password
	hashed, _, _ := DefaultHash(password)
	u.Retrieve() // Lookup from the database
	return hashed == u.PasswordHash
}

// ValidateAuthToken verifies the given auth token is valid for the user.
func (u *UserAuth) ValidateAuthToken(token string) bool {
	u.Retrieve() // Lookup from the database
	return u.Token == token
}

/*
	DBActor
*/

// Create inserts a new UserAuth into the database
func (u *UserAuth) Create(conn *sql.DB) error {
	_, err := conn.Exec(`
	INSERT INTO auth.UserAuth (
		id,
		username,
		account_creation,
		enabled,
		password_hash,
		password_salt,
		timestamp)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		u.ID, u.Username, u.FirstName, u.LastName, u.Email, u.Enabled, u.PasswordHash, u.PasswordSalt, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Retrieve a UserAuth from the DB by filtering on Username
func (u *UserAuth) Retrieve(conn *sql.DB) error {
	err := conn.QueryRow(`
	SELECT (
		id,
		username,
		account_creation,
		enabled,
		password_hash,
		password_salt,
		timestamp)
	FROM auth.UserAuth
	WHERE username=$1`,
		u.Username).Scan(u)
	if err != nil {
		return err
	}
	return nil
}

// Update a user entry in the database
func (u *UserAuth) Update(conn *sql.DB) error {
	_, err := conn.Exec(`
	UPDATE auth.UserAuth
	SET first_name='$2', last_name='$3', email='$4', enabled=$5, timestamp=$6
	WHERE username=$1`,
		u.Username, u.Enabled, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a UserAuth record from the database. In most cases it will
// probably be best practice to simply flag a user as disabled via a PUT, but
// we do also need to expose this ability.
func (u *UserAuth) Delete(conn *sql.DB) error {
	err := conn.QueryRow(`
	DELETE FROM auth.UserAuth
	WHERE username=$1`, u.Username).Scan(u)
	if err != nil {
		return nil, err
	}
}

/*
	RESTfulHandler
*/

// HandlePost creates a new UserAuth record.
func (u *UserAuth) HandlePost(w http.ResponseWriter, req *http.Request) {
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

// HandleGet returns the specified UserAuth record
func (u *UserAuth) HandleGet(w http.ResponseWriter, req *http.Request) {
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

// HandlePut updates a UserAuth resource.
func (u *UserAuth) HandlePut(w http.ResponseWriter, req *http.Request) {
	// Parse body of PUT request into a UserAuth struct
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
func (u *UserAuth) HandleDelete(w http.ResponseWriter, req *http.Request) {
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
func (u *UserAuth) HTTPPost(serverURL string) (resp *http.Response, err error) {
	endpoint := "/users"
	return torque.PostJSON(endpoint, u)
}

// HTTPGet requests a user record from the server.
func (u *UserAuth) HTTPGet(serverURL string) (resp *http.Response, err error) {
	return nil, nil
}

// HTTPPut updates the server with the current state of the UserAuth record.
func (u *UserAuth) HTTPPut(serverURL string) (resp *http.Response, err error) {
	return nil, nil
}

// HTTPDelete deletes the matching user record on the server.
func (u *UserAuth) HTTPDelete(serverURL string) (resp *http.Response, err error) {
	return nil, nil
}
