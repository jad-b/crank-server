package redteam

import (
	"github.com/jad-b/torque"
	"github.com/jad-b/torque/users"
)

// GetOrCreateUser looks for a User in the database. If none is found, it
// creates a user with the requested username & password.
func GetOrCreateUser(username, password string) (*users.UserAuth, error) {
	u := &users.UserAuth{Username: username}
	if err := u.Retrieve(torque.DB); err != nil {
		return u, err
	}
	// Create a user instead
	u.PasswordHash = password
	if err := u.Create(torque.DB); err != nil {
		return u, err
	}
	return u, nil
}
