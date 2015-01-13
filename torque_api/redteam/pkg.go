package redteam

import (
	"flag"
	"log"

	"github.com/jad-b/torque"
)

const (
	username = "redteam"
	password = "redteam"
)

var (
	https      = flag.Bool("secure", false, "Whether to use HTTPS")
	torqueAddr = torque.HostPortFlag{Host: "localhost:18000"}
)

func init() {
	flag.Var(&torqueAddr, "torque-addr", "host:port of Torque API")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
