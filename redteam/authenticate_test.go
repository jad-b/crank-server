package redteam

import (
	"flag"
	"testing"
)

/*
	authenticate logs into a torque server
*/

const (
	username = "redteam"
	password = "redteam"
)

var (
	torqueURL = flag.String("torque-url", "localhost:8080", "URL of Torque API")
)

func TestAuthentication(t *testing.T) {
	_, err := AuthenticateToServer(*torqueURL, username, password)
	if err != nil {
		t.Error(err)
	}
}
