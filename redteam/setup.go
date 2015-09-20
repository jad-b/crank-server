package redteam

import (
	"errors"
	"log"
	"strings"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/client"
	"github.com/jad-b/torque/users"
	"github.com/jmoiron/sqlx"
)

// CreateAdminUser inserts an admin user into the database.
func CreateAdminUser(db *sqlx.DB, adminUser, adminPassword string) error {
	log.Printf("Creating admin account for %s", adminUser)
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
	log.Print("Creating redteam user")
	// Connect to the server
	tAPI := client.NewTorqueAPI(server)

	// Authenticate our user
	err := tAPI.Authenticate(adminUser, adminPassword)
	if err != nil {
		return err
	}
	log.Print("Admin user authenticated")

	// Create redteam user account
	redUser := users.UserAuth{
		Username:     username,
		PasswordHash: password,
	}
	resp, err := tAPI.Post(&redUser)
	if err != nil {
		// Not a problem if the user already exists
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil
		}
		return err
	}
	// Unsuccessful creation through the API
	if resp.StatusCode != 200 {
		torque.LogResponse(resp)
		// Parse error from resp
		var errResp torque.ErrorResponse
		if err != nil {
			return errors.New("Failed to read API error")
		}
		return errResp
	}
	log.Print("Redteam user created")
	// Overwrite redteam user with returned fields
	return nil
}
