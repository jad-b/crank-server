package users

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jad-b/torque"
)

// UserSecrets is a struct encapsulating a user's sensitive data. This has been
// extracted from the user tables, in order to isolate sensitive data from
// non-sensitive data so we can reduce our risk of doing something stupid down
// the road.
//
// SQL:
// CREATE TABLE ppi (
//	 "id" numeric(5,2) NOT NULL
//	 "user_id" text NOT NULL UNIQUE
//	 "password_hash" text NOT NULL
//	 "password_salt" text NOT NULL
//	 "iteration_count" numeric(5,2) NOT NULL
//   "current_token" text NOT NULL
//	 "timestamp" timestamp(0) with time zone NOT NULL UNIQUE,
//   "token_last_seen" timestamp(0) with time zone NOT NULL UNIQUE,
// );
type UserSecrets struct {
	ID             int `json:"id"`
	UserID         int `json:"user_id"`
	PasswordHash   string
	PasswordSalt   string
	IterationCount int
	Timestamp      time.Time `json:"timestamp"`
	CurrentToken   string    `json:"token"`
	TokenLastSeen  time.Time `json:"token_last_seen"`
}

// Create a new UserSecrets instance
func NewUserSecrets(u *User, passwordHash string) *UserSecrets {
	return &UserSecrets{
		UserID:         u.ID,
		PasswordHash:   passwordHash,
		PasswordSalt:   NewSalt(DefaultSaltLength),
		IterationCount: DefaultIterationCount,
		Timestamp:      time.Now(),
	}
}

func (p *UserSecrets) Validate(password string) bool {}
func GenerateHash(password string) string            {}

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
)

// Generate a new, random, salt of the length specified
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
