package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/metrics"
)

var (
	web      bool
	registry = map[string]torque.CommandLineActor{
		"bodyweight": &metrics.Bodyweight{},
	}
	addr    = flag.String("addr", "", "Host:port of Torque server")
	verbose = flag.Bool("v", false, "Toggle verbose output")
)

/* cli is the command-line interface for Torque.

Command syntax:
	torque <options> <action> <resource> arguments>
*/
func main() {
	log.SetOutput(os.Stderr)
	flag.Parse()
	dbOrWeb()

	// Handle all errors generically
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s is an invalid call of torque", os.Args)
		}
	}()

	// pass the remaining args off to the resources to handle
	handleArgs()
}

// dbOrWeb determines whether we're talking HTTP to a web server or directly to
// a database.
func dbOrWeb() {
	if *addr != "" {
		web = true
	} else {
		web = false
		// Open up a Database connection
		pgconf := torque.LoadPostgresConfig()
		torque.OpenDBConnection(pgconf)
	}
}

func handleArgs() {
	// Check we received a minimal amount of arguments
	remainder := flag.Args()
	lenRemainder := len(remainder)
	// log.Printf("Remaining args: %s", remainder)
	if lenRemainder < 1 {
		log.Panic("No resource specified")
	} else if lenRemainder < 2 {
		log.Panic("No action specified")
	} else if lenRemainder < 3 {
		log.Panic("No data was provided")
	}

	// Delegate remaining arg parsing to the identified resource
	resource, action := remainder[0], remainder[1]
	r, ok := registry[resource]
	if !ok {
		log.Fatalf("%s not recognized as resource", remainder[1])
	}
	// Resource located; have it parse the remaining flags.
	r.ParseFlags(action, remainder[2:])
	// Determine if we're going over HTTP or directly to the database
	var err error
	if web {
		var resp *http.Response
		resp, err = torque.ActOnWebServer(r, action, *addr)
		if *verbose {
			var buf bytes.Buffer
			resp.Write(&buf)
			log.Printf("Response from %s:\n%s \n", *addr, buf.String())
		}
	} else {
		err = torque.ActOnDB(r, action, torque.DBConn)
	}
	if err != nil {
		log.Panic(err)
	}
}
