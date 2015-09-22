package users

import (
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/jad-b/torque"
	"github.com/jmoiron/sqlx"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultSaltLength is the default length for newly generated PasswordSalts
	DefaultSaltLength = 32

	// DefaultBcryptCost is the power-of-two iterations (2^cost) to apply via
	// bcrypt when hashing.
	DefaultBcryptCost = 12

	// AuthTokenLength is the size of a generated auth token
	AuthTokenLength = 32
)

var (
	// AuthTokenLifespan is how long an AuthToken should be valid for from its
	// creation time.
	AuthTokenLifespan = time.Hour * 1
)

// HandleAuthentication validates username & password and returns a User object
// with a new AuthToken, or an Unauthorized error.
func HandleAuthentication(w http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()
	if !ok {
		w.Header().Set(torque.HeaderAuthenticate, "Your bad")
		http.Error(w, "Failed to retrieve credentials from request", http.StatusUnauthorized)
		return
	}
	log.Printf("Authentication request from %s", username)
	// Check the credentials
	user := UserAuth{Username: username}
	if err := user.Retrieve(torque.DB); err != nil {
		e := torque.ErrorResponse{"Invalid credentials"}
		w.Header().Set(torque.HeaderAuthenticate, e.Error())
		torque.HTTPError(w, e, http.StatusUnauthorized)
		return
	}
	ok = user.ValidatePassword(password)
	if !ok {
		e := torque.ErrorResponse{"Invalid credentials"}
		w.Header().Set(torque.HeaderAuthenticate, e.Error())
		torque.HTTPError(w, e, http.StatusUnauthorized)
		return
	}
	log.Print("Authenticated ", username)
	// Assign user an auth token
	user.Authorize(torque.DB)
	if err := user.Update(torque.DB); err != nil {
		e := torque.ErrorResponse{"Failed to issue authorization token"}
		w.Header().Set(torque.HeaderAuthenticate, e.Error())
		torque.HTTPError(w, e, http.StatusInternalServerError)
		return
	}
	log.Print("Authorized ", username)
	// Set Authorization header
	w.Header().Set(AuthHeader(&user))
	// Send user object back with our request
	torque.WriteOkayJSON(w, user)
}

// BuildAuthenticationRequest prepares a HTTP request for retrieving an auth
// token.
func BuildAuthenticationRequest(server, username, password string) (*http.Request, error) {
	// Prepare the URL
	u := url.URL{Scheme: torque.Scheme, Host: server}
	// Append the Authentication path to our URL
	u.Path = torque.SlashJoin(u.Path, "authenticate/")
	// Create the HTTP request
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	return req, nil
}

// AuthorizeAs determines if A is allowed to act on B.
// It returns the authorized user ID. This relies on resources supplying who
// they *think* should be the owner.
func AuthorizeAs(db *sqlx.DB, token string, owner int) (int, error) {
	// Retrieve auth'd user account
	actor, err := SwapTokenForUser(db, token)
	if err != nil {
		return 0, err
	}
	// If no resource owner was provided, it defaults to whoever owns the auth
	// token
	if owner == 0 {
		return actor.ID, nil
	}
	// Obviously users are allowed to be themselves
	if actor.ID == owner {
		return owner, nil
	}
	// Superusers are allowed to impersonate other users
	if actor.Superuser {
		log.Printf("%s is a superuser; proceeding", actor.Username)
		return owner, nil
	}
	return 0, errors.New("Unauthorized user")
}

// SwapTokenForUser retrieves the user using the issued auth token.
func SwapTokenForUser(db *sqlx.DB, token string) (*UserAuth, error) {
	// Retrieve actor's account
	var user UserAuth
	err := db.Get(
		&user,
		fmt.Sprintf(`
			SELECT
				id,
				username,
				superuser
			FROM %s.%s
			WHERE current_token=$1`,
			Schema,
			userAuthTableName),
		token)
	if err != nil { // Token didn't exist
		return nil, err
	}
	return &user, nil
}

// SwapTokenForID retrieves the user using the issued auth token.
func SwapTokenForID(db *sqlx.DB, token string) (int, error) {
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
	if err != nil { // Token didn't exist
		return 0, err
	}
	return id, nil
}

// AuthHeader builds the Authorization header.
func AuthHeader(u *UserAuth) (string, string) {
	return torque.HeaderAuthorization, fmt.Sprintf("token token=%s", u.CurrentToken)
}

// ParseAuthToken extracts the auth token from a string, presumably the
// Authorization header.
func ParseAuthToken(authHeader string) (token string) {
	// everything after `token=` and before a comma
	re := regexp.MustCompile(`token=(?P<token>.+)`)
	matches := re.FindStringSubmatch(authHeader)
	return matches[1]
}

// DefaultHash applies a one-way bcrypt hash to a string.
// It returns the resulting hash, the salt used, and the cost (power of two of
// iterations to be performed). Good for creating passwords.
func DefaultHash(password string) (hash, salt string, cost int) {
	s, err := GenerateRandomString(DefaultSaltLength)
	if err != nil {
		log.Panic(err)
	}
	return GenerateHash(password, s, DefaultBcryptCost), s, DefaultBcryptCost
}

// GenerateHash one-way bcrypt hashes the password.
// Good for verifying passwords on existing user accounts.
// TODO(jdb) Salt is currently unused
func GenerateHash(password, salt string, cost int) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Panic(err)
	}
	return string(hashed)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// From: https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// From: https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
