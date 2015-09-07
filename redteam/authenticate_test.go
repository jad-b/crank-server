// +build test rest

package redteam

import (
	"flag"
	"log"
	"testing"

	"github.com/jad-b/torque"
)

/*
	authenticate logs into a torque server
*/

const (
	username = "redteam"
	password = "redteam"
)

var (
	torqueAddr = torque.HostPortFlag{Host: "localhost:18000"}
	https      = flag.Bool("secure", false, "Whether to use HTTPS")
)

func init() {
	flag.Var(&torqueAddr, "torque-addr", "host:port of Torque API")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestAuthentication(t *testing.T) {
	if *https {
		torqueAddr.Scheme = "https"
	} else {
		torqueAddr.Scheme = "http"
	}
	t.Log("Authenticating against ", torqueAddr.String())
	_, err := AuthenticateToServer(torqueAddr.String(), username, password)
	if err != nil {
		t.Error(err)
	}
}
