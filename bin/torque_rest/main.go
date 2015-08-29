package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/metrics"
)

var (
	addr = flag.String("addr", "127.0.0.1:8000", "Host:port of Torque server")
	cert = flag.String("cert", "", "TLS server certificate")
	key  = flag.String("key", "", "TLS server private key")
)

func runServer() {
	log.Print("Starting server")
	defer log.Fatal("Stopping server")

	// Register RESTfulHandlers
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics/bodyweight/",
		torque.RouteRequest(&metrics.Bodyweight{}))

	// Setup our database connection
	pgConf := torque.LoadPostgresConfig()
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
	log.SetOutput(os.Stderr)

	runServer()
}
