package main

import (
	"flag"
	"log"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/redteam"
)

var (
	adminUser     = flag.String("admin-user", "", "Admin account username")
	adminPassword = flag.String("admin-password", "", "Admin account password")
	torqueAddr    = torque.HostPortFlag{Host: "localhost:18000"}
	// Modal flags
	setupMode   = flag.Bool("setup", false, "Setup mode")
	redteamMode = flag.Bool("redteam", false, "Redteam mode")
	adminMode   = flag.Bool("admin", false, "Admin mode")
)

func main() {
	flag.Var(&torqueAddr, "localhost:18000", "Address of Torque server")
	flag.Parse()

	// Flag switching
	if *setupMode {
		if *adminUser == "" || *adminPassword == "" {
			log.Fatal("Missing admin credential(s)")
			return
		}
		var err error
		if *adminMode {
			db := torque.Connect()
			defer db.Close()
			err = redteam.CreateAdminUser(db, *adminUser, *adminPassword)
		}
		if *redteamMode {
			err = redteam.SetupRedteamUser(*adminUser, *adminPassword, torqueAddr.String())
		}
		if err != nil {
			log.Print(err)
		}
	}
}
