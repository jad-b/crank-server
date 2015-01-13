package torque

import "log"

// Torque-wide initialization
func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
