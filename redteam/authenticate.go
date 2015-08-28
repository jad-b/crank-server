package redteam

import "github.com/jad-b/torque/client"

/*
	authenticate logs into a torque server
*/

// AuthenticateToServer goes through the authentication workflow
func AuthenticateToServer(serverURL, username, password string) (*client.TorqueAPI, error) {
	// Connect to the server
	tAPI := client.NewTorqueAPI(serverURL)

	// Authenticate our user
	err := tAPI.Authenticate(username, password)
	return tAPI, err
}
