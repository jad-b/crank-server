package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jad-b/torque/metrics"
)

/* cli is the command-line interface for Torque.

Command syntax:
	torque <options> <action> <resource> arguments>
*/

func handleArgs() {
	// Check we received a minimal amount of arguments
	remainder := flag.Args()
	lenRemainder := len(remainder)
	log.Printf("Remaining args: %s", remainder)
	if lenRemainder < 1 {
		log.Printf("No action specified")
	} else if lenRemainder < 2 {
		log.Printf("No action specified")
	} else if lenRemainder < 3 {
		log.Printf("No data was provided")
	}

	// Delegate remaining arg parsing to the identified resource
	resource, action := remainder[0], remainder[1]
	switch resource {
	case "bodyweight":
		bw := &metrics.Bodyweight{}
		ret := bw.ParseFlags(action, remainder[2:])
		fmt.Printf("%s %+v\n", action, ret)
	default:
		fmt.Printf("%s not recognized as resource", remainder[1])
	}
}

func main() {
	flag.Parse()
	log.SetOutput(os.Stderr)

	// Handle all errors generically
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s is an invalid call of torque", os.Args)
		}
	}()

	handleArgs()
}
