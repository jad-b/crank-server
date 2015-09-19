package redteam

import (
	"github.com/jad-b/torque"
	"github.com/jad-b/torque/client"
	"github.com/jad-b/torque/users"
)

//

// SetupRedteamUser performs necessary tasks to operate the redteam user
func SetupRedteamUser(adminUser, adminPassword, server string) error {
	// Create acount on Torque
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
