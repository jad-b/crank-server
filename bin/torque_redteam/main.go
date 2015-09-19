package main

import (
	"flag"
	"log"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/redteam"
)

var (
	admin         = flag.String("admin", "", "Admin account username")
	adminPassword = flag.String("admin-pasword", "", "Admin account password")
	torqueAddr    = torque.HostPortFlag{Host: "localhost:18000"}
)

func main() {
	flag.Var(&torqueAddr, "localhost:18000", "Address of Torque server")
	flag.Parse()

	err := redteam.SetupRedteamUser(*admin, *adminPassword, torqueAddr.String())
	log.Print(err)
}
