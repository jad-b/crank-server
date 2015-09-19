package redteam

import (
	"strings"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/client"
	"github.com/jad-b/torque/users"
	"github.com/jmoiron/sqlx"
)

// CreateAdminUser inserts an admin user into the database.
func CreateAdminUser(db *sqlx.DB, adminUser, adminPassword string) error {
	// Create admin user account
	admin := users.NewUserAccount(adminUser, adminPassword)
	admin.Superuser = true
	// Insert into the database
	if err := admin.Create(db); err != nil {
		// Not a problem if the user already exists
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil
		}
		return err
	}
	return nil
}

// SetupRedteamUser performs necessary tasks to operate the redteam user
func SetupRedteamUser(adminUser, adminPassword, server string) error {
	// Connect to the server
	tAPI := client.NewTorqueAPI(server)

	// Authenticate our user
	err := tAPI.Authenticate(adminUser, adminPassword)

	// Create redteam user account
	redUser := users.NewUserAccount(username, password)
	resp, err := tAPI.Post(redUser)
	// TODO Handle account already existing
	if err != nil {
		return err
	}
	// Overwrite redteam user with returned fields
	err = torque.ReadJSONResponse(resp, &redUser)
	if err != nil {
		return err
	}
	return nil
}
