package users

import (
	crand "crypto/rand"
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/jad-b/torque"

	"golang.org/x/crypto/bcrypt"
)

const (
	// The alphabet constant is used to expose all of the valid characters that
	// we can use when generating a new PasswordSalt. This constant can be
	// safelty updated without running the risk of breaking already generated
	// salts. Because salts are stored on a per-user basis and are only
	// generated/replaced when a user first creates their account or creates a
	// new password.
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// DefaultSaltLength is the default length for newly generated PasswordSalts
	DefaultSaltLength = 32

	// DefaultIterationCount is the default number of times to iterate over the
	// password while hashing. Using a higher iteration count will increase the
	// cost of an exhaustive search but will also make derivation
	// proportionally slower. We'll likely need to fine tune the default
	// iteration count as time goes on to provide a better user experience,
	// while still maintaining a higher level of security. We'll store the
	// iteration count along side each salt and hash in order to allow us the
	// flexibility to safetly modify this default in the future.
	DefaultIterationCount = 1000
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
		w.Header().Set("WWW-Authenticate", "Your bad")
		http.Error(w, "Failed to retrieve credentials from request", http.StatusUnauthorized)
		return
	}
	log.Printf("Authentication request from %s", username)
	// Check the credentials
	user := UserAuth{Username: username}
	ok = user.ValidatePassword(password)
	if !ok {
		e := torque.ErrorResponse{"Invalid credentials"}
		w.Header().Set("WWW-Authenticate", e.Error())
		torque.HTTPError(w, e, http.StatusUnauthorized)
		return
	}
	// Create an auth token

}

// DefaultHash applies a one-way bcrypt hash to a string.
// It returns the resulting hash, the salt used, and the cost (power of two of
// iterations to be performed). Good for creating passwords.
func DefaultHash(password string) (hash, salt string, cost int) {
	s := NewSalt(DefaultSaltLength)
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

// NewSalt generates a new, random, salt of the length specified
func NewSalt(length int) string {
	// Create a new []byte of size *length*
	b := make([]byte, length)

	// For each entry in our new []byte, get a random integer within the range
	// of our constant alphabet, and insert alphabet[random_int] into our new
	// byte array
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
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
