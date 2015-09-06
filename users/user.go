package users

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

const (
	userAuthTableName = "user_auth"
)

var (
	// UserAuthTable describes the SQL fields
	userAuthTable = `
id serial PRIMARY KEY,
username text NOT NULL UNIQUE,
account_created timestamp(0) WITH time zone NOT NULL,
enabled boolean DEFAULT TRUE,
superuser boolean DEFAULT FALSE,
password_hash text NOT NULL,
password_salt text NOT NULL,
cost integer NOT NULL,
current_token text,
token_created timestamp(0) WITH time zone,
token_last_used timestamp(0) WITH time zone`
)

// CreateTableUserAuth creates the UserAuth table
func CreateTableUserAuth(db *sqlx.DB) error {
	err := torque.CreateTable(
		db,
		Schema,
		userAuthTableName,
		userAuthTable,
		true)
	return err
}

// UserAuth holds the auth token data mapped to a UserAuth ID.
// TODO It also holds some data, like Enabled, AccountCreated, that might be better
// served in a UserMeta table, but I left it here for simplicity's sake. Also, I
// hear JOINs are expensive.
type UserAuth struct {
	// Primary key. Links user tables together
	ID int `json:"id"`
	// Keep username close to passwordHash for authentication calls
	Username       string    `json:"username"`
	AccountCreated time.Time `json:"account_created" db:"account_created"`
	Enabled        bool      `json:"enabled"`
	Superuser      bool      `json:"superuser"`
	// Salt used to hash the password
	PasswordSalt string `json:"-" db:"password_salt"`
	// Type of hash used
	PasswordHash string `json:"-" db:"password_hash"`

	// Power-of-two times we iterated over the stored password when hashing
	Cost int `json:"-"`

	// Currently active auth token
	CurrentToken string    `json:"current_token" db:"current_token"`
	TokenCreated time.Time `json:"token_created" db:"token_created"`
	// Last time the token was used in an API request
	TokenLastUsed time.Time `json:"token_last_used" db:"token_last_used"`
}

// IsAuthenticated returns whether or not the User is authenticated. This
// method *does not* touch the server, so it relies upon available client-side
// data.
func (u *UserAuth) IsAuthenticated() bool {
	return u.CurrentToken != "" &&
		time.Since(u.TokenCreated) < AuthTokenLifespan
}

// NewUserAccount creates a new UserAuth instance with some defaults in place.
func NewUserAccount(username, password string) *UserAuth {
	hash, salt, cost := DefaultHash(password)
	return &UserAuth{
		Username:       username,
		Enabled:        true,
		PasswordHash:   hash,
		PasswordSalt:   salt,
		Cost:           cost,
		AccountCreated: time.Now(),
	}
}

// Authorize creates a new auth token and updates its metadata fields.
// It expects that the previous User record has been fully loaded from the
// database, else you risk overwriting an existing record!
func (u *UserAuth) Authorize(db *sqlx.DB) error {
	token, err := GenerateRandomString(AuthTokenLength)
	if err != nil {
		return err
	}
	u.CurrentToken = token
	now := time.Now()
	u.TokenCreated = now
	u.TokenLastUsed = now
	// Save changes to DB
	u.Update(db)
	return nil
}

// ValidatePassword verifies if the given username/password is valid.
func (u *UserAuth) ValidatePassword(password string) bool {
	err := u.Retrieve(torque.DB) // Lookup from the database
	if err != nil {              // User not found
		log.Printf("User %s not found", u.Username)
		return false
	}
	// Hash the password
	hashed, _, _ := DefaultHash(password)
	ok := hashed == u.PasswordHash
	if !ok {
		log.Print("Invalid login for ", u.Username)
	}
	return ok
}

// ValidateAuthToken verifies the given auth token is valid for the user.
func (u *UserAuth) ValidateAuthToken(token string) bool {
	err := u.Retrieve(torque.DB) // Lookup from the database
	if err != nil {              // User not found
		log.Printf("User %s not found", u.Username)
		return false
	}

	ok := token == u.CurrentToken
	if !ok {
		log.Print("Invalid token for ", u.Username)
	}
	return ok
}

/*
	DBActor
*/

// Create inserts a new UserAuth row into the database
func (u *UserAuth) Create(db *sqlx.DB) error {
	_, err := db.NamedExec(fmt.Sprintf(`
	INSERT INTO %s.%s (
		username,
		account_created,
		enabled,
		superuser,
		password_hash,
		password_salt,
		cost,
		current_token,
		token_created,
		token_last_used
	) VALUES (
		:username,
		:account_created,
		:enabled,
		:superuser,
		:password_hash,
		:password_salt,
		:cost,
		:current_token,
		:token_created,
		:token_last_used
	)`, Schema, userAuthTableName), u)
	return err
}

// Retrieve a UserAuth from the DB by filtering on Username
func (u *UserAuth) Retrieve(db *sqlx.DB) error {
	return db.Get(
		u,
		fmt.Sprintf(`
			SELECT *
			FROM %s.%s
			WHERE username=$1`,
			Schema,
			userAuthTableName),
		u.Username)
}

// Update a user entry in the database
// TODO(jdb) Might be a bad idea to override everything - kind of implies you'l
// want to RETRIEVE the existing record, apply changes, then UPDATE the row.
// Or maybe we should implement PATCH for partial updates.
func (u *UserAuth) Update(db *sqlx.DB) error {
	_, err := db.NamedExec(
		fmt.Sprintf(`
			UPDATE %s.%s
			SET
				account_created=:account_created,
				enabled=:enabled,
				superuser=:superuser,
				password_hash=:password_hash,
				password_salt=:password_salt,
				cost=:cost,
				current_token=:current_token,
				token_created=:token_created,
				token_last_used=:token_last_used
			WHERE username=:username`,
			Schema,
			userAuthTableName),
		u)
	return err
}

// Delete removes a UserAuth record from the database. In most cases it will
// probably be best practice to simply flag a user as disabled via a PUT, but
// we do also need to expose this ability.
func (u *UserAuth) Delete(db *sqlx.DB) error {
	_, err := db.NamedExec(
		fmt.Sprintf(`
			DELETE FROM %s.%s
			WHERE username=:username`,
			Schema,
			userAuthTableName),
		u)
	return err
}

/*
	RESTfulHandler
*/

// HandlePost creates a new UserAuth record.
func (u *UserAuth) HandlePost(w http.ResponseWriter, req *http.Request) {
	// Retrieve username & password from Basic-Auth header
	username, password, ok := req.BasicAuth()
	if !ok {
		http.Error(w, "No username and password provided for account creation",
			http.StatusBadRequest)
		return
	}
	// Create user account
	u = NewUserAccount(username, password)
	// Save to database
	if err := u.Create(torque.DB); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	torque.WriteOkayJSON(w, u)
}

// HandleGet returns the specified UserAuth record
func (u *UserAuth) HandleGet(w http.ResponseWriter, req *http.Request) {
	userID, err := parseUserID(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.ID = userID
	if err = u.Retrieve(torque.DB); err != nil {
		http.NotFound(w, req)
		return
	}
	torque.WriteOkayJSON(w, u)
}

// HandlePut updates a UserAuth resource.
func (u *UserAuth) HandlePut(w http.ResponseWriter, req *http.Request) {
	userID, err := parseUserID(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.ID = userID
	// Parse body of PUT request into a UserAuth struct
	err = torque.ReadBodyTo(w, req, u)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = u.Update(torque.DB); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	// Write updated record to client
	torque.WriteOkayJSON(w, u)
}

// HandleDelete removes a user record from the database.
func (u *UserAuth) HandleDelete(w http.ResponseWriter, req *http.Request) {
	userID, err := parseUserID(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.ID = userID
	if err = u.Delete(torque.DB); err != nil {
		http.NotFound(w, req)
		return
	}
	torque.WriteOkayJSON(w, u)
}

func parseUserID(earl *url.URL) (int, error) {
	// Parse user ID from URL
	re := regexp.MustCompile(`/users/(\d+)/?`)
	parts := re.FindStringSubmatch(earl.Path)
	if len(parts) < 2 {
		return 0, fmt.Errorf("Expected /users/{user_id}/, not %s", earl.Path)
	}
	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("%s is an invalid user ID", parts[1])
	}
	return userID, nil
}

/*
	RESTfulResource
*/

// GetResourceName returns the name UserAuth wishes to be referred to by in the
// URL
func (u *UserAuth) GetResourceName() string {
	return torque.SlashJoin("users", u.Username)
}
