package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/metrics"
	"github.com/jad-b/torque/users"
)

var (
	addr = flag.String("addr", "127.0.0.1:18000", "Host:port of Torque server")
	cert = flag.String("cert", "", "TLS server certificate")
	key  = flag.String("key", "", "TLS server private key")
)

func runServer() {
	log.Printf("Serving on %s", *addr)
	defer log.Fatal("Stopping server")

	// Register RESTfulHandlers
	mux := http.NewServeMux()
	bw := &metrics.Bodyweight{}
	mux.HandleFunc("/"+bw.GetResourceName(), torque.RouteRequest(bw))
	mux.Handle("/authenticate/", torque.SmartHandler(users.HandleAuthentication))
	mux.HandleFunc("/users/", torque.RouteRequest(&users.UserAuth{}))
	// Default handler - do nothing
	mux.HandleFunc("/", torque.LogRequestThenError)

	// Setup our database connection
	pgConf := torque.LoadPostgresConfig(*torque.PsqlConf)
	torque.OpenDBConnection(pgConf)

	// Start the server
	if *cert != "" && *key != "" {
		http.ListenAndServeTLS(*addr, *cert, *key, mux)
	} else {
		http.ListenAndServe(*addr, mux)
	}
}

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)

	runServer()
}
