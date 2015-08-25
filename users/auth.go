package users

import "math/rand"

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
)

// DefaultHash applies a one-way bcrypt hash to a string.
// It returns the resulting hash, the salt used, and the cost (power of two of
// iterations to be performed).
func DefaultHash(password string) (hash, salt string, cost int) {
	s = NewSalt(DefaultSaltLength)
	return GenerateHash(password, s, DefaultBcryptCost), s, DefaultBcryptCost
}

// GenerateHash one-way bcrypt hashes the password.
// TODO(jdb) Salt is currently unused
func GenerateHash(password, salt string, cost int) string {
	hashed, err := bcrypt.GeneratePasswordFrom(password, cost)
	if err != nil {
		log.Error(err)
	}
	return hashed
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
