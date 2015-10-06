package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"
)

const (
	userAuthTableName = "user_auth"
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
	return torque.CreateTable(
		db,
		Schema,
		userAuthTableName,
		userAuthTable,
		true)
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
	PasswordHash string `json:"password_hash" db:"password_hash"`
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
	log.Printf("Issused auth token for user '%s'", u.Username)
	return nil
}

// ValidatePassword verifies if the given username/password is valid.
func (u *UserAuth) ValidatePassword(password string) bool {
	log.Printf("Validating %s/%s", u.Username, password)
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

// ValidateAuthToken verifies the given auth token is valid for the user.
func (u *UserAuth) ValidateAuthToken(token string) bool {
	return token == u.CurrentToken
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

// GetUserByToken retrieves the UserAuth row using the auth token
func GetUserByToken(db *sqlx.DB, token string) (*UserAuth, error) {
	var user UserAuth
	err := db.Get(
		&user,
		fmt.Sprintf(`
			SELECT *
			FROM %s.%s
			WHERE current_token=$1`,
			Schema,
			userAuthTableName),
		token)
	return &user, err
}

// GetUserIDByToken does a User ID lookup using the Auth Token
func GetUserIDByToken(db *sqlx.DB, token string) int {
	var id int
	err := db.Get(
		&id,
		fmt.Sprintf(`
			SELECT id
			FROM %s.%s
			WHERE current_token=$1`,
			Schema,
			userAuthTableName),
		token)
	if err != nil {
		log.Print(err)
		return 0
	}
	return id
}

/*
	RESTfulHandler
*/

// HandlePost creates a new UserAuth record.
// It's a bit of a chicken & egg problem, in that you must be an authorized user to create users.
// Otherwise, anyone could just blast the API with user creations.
func (u *UserAuth) HandlePost(w http.ResponseWriter, req *http.Request) {
	log.Print("Request: create user")
	// Get UserID from token
	token := ParseAuthToken(req.Header.Get(torque.HeaderAuthorization))
	_, err := SwapTokenForID(torque.DB, token)
	if err != nil {
		torque.HTTPError(w, errors.New("User not authorized"),
			http.StatusUnauthorized)
		return
	}
	// Extract user to create from request body
	var userBody UserAuth
	err = torque.ReadJSONRequest(req, &userBody)
	if err != nil {
		torque.HTTPError(w, errors.New("User not found in request body"),
			http.StatusBadRequest)
		return
	}
	// Setup user account
	newUser := NewUserAccount(userBody.Username, userBody.PasswordHash)
	if err := CheckPasswordStrength(userBody.PasswordHash); err != nil {
		log.Printf("Password sucked: %s", userBody.PasswordHash)
		torque.HTTPError(w, err, http.StatusBadRequest)
		return
	}
	log.Printf("Credentials are present, creating %s/%s", userBody.Username, userBody.PasswordHash)
	// Save to database
	if err := newUser.Create(torque.DB); err != nil {
		torque.HTTPError(
			w,
			fmt.Errorf("Failed to create user account %s in database", u.Username),
			http.StatusInternalServerError)
		return
	}
	torque.WriteOkayJSON(w, newUser)
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
	err = torque.ReadJSONRequest(req, u)
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
	return Category + "/"
}
