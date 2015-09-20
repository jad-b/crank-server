// +build test rest

package redteam

import "testing"

/*
	authenticate logs into a torque server
*/

func TestAuthentication(t *testing.T) {
	if *https {
		torqueAddr.Scheme = "https"
	} else {
		torqueAddr.Scheme = torque.Scheme
	}
	t.Log("Authenticating against ", torqueAddr.String())
	_, err := AuthenticateToServer(torqueAddr.String(), username, password)
	if err != nil {
		t.Error(err)
	}
}
