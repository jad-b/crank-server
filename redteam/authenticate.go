package redteam

import "github.com/jad-b/torque"

/*
	authenticate logs into a torque server
*/

// AuthenticateToServer goes through the authentication workflow
func AuthenticateToServer(serverURL, username, password string) (torque.API, torque.UserAuth) {
	// Connect to the server
	sURL := "http://localhost:8080"
	tAPI := torque.NewTorqueAPI(sURL)

	// Authenticate our user
	u := tAPI.authenticate(username, password)
	return tAPI, u
}
