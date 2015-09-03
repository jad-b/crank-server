package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
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
	// The error that killed the program. Having this as a script-global allows
	// us to set an error wherever, recover generically with a 'defer' in
	// main(), and still output something meaningful.
	terminalError error
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
			fmt.Printf("%s is an invalid call of torque: %s\n", os.Args, terminalError)
			flag.Usage()
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
		pgconf := torque.LoadPostgresConfig(*torque.PsqlConf)
		torque.OpenDBConnection(pgconf)
	}
}

func handleArgs() {
	// Check we received a minimal amount of arguments
	remainder := flag.Args()
	lenRemainder := len(remainder)
	// log.Printf("Remaining args: %s", remainder)
	if lenRemainder < 1 {
		terminalError = errors.New("No resource specified")
	} else if lenRemainder < 2 {
		terminalError = errors.New("No action specified")
	} else if lenRemainder < 3 {
		terminalError = errors.New("No data was provided")
	}

	// Delegate remaining arg parsing to the identified resource
	resource, action := remainder[0], remainder[1]
	r, ok := registry[resource]
	if !ok {
		terminalError = fmt.Errorf("%s not recognized as resource", remainder[1])
	}
	// Resource located; have it parse the remaining flags.
	r.ParseFlags(action, remainder[2:])
	// Determine if we're going over HTTP or directly to the database
	var err error
	if web {
		var resp *http.Response
		resp, err = torque.ActOnWebServer(action, *addr)
		if *verbose {
			var buf bytes.Buffer
			resp.Write(&buf)
			log.Printf("Response from %s:\n%s \n", *addr, buf.String())
		}
	} else {
		err = torque.ActOnDB(r, action, torque.DB)
	}
	if err != nil {
		terminalError = err
	}
}
