package api

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jad-b/torque/ui"
)

var (
	addr = flag.String("addr", "127.0.0.1", "IP address to bind server to")
	port = flag.String("port", "8000", "Port to bind server to")
	cert = flag.String("cert", "", "TLS server certificate")
	key  = flag.String("key", "", "TLS server private key")
)

// IdentityHandler echoes the hostname back to the client
func IdentityHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("host is %s", req.Host)
	fmt.Fprintf(w, "%s, this is me.", req.Host)
}

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)
	log.Print("Starting server")
	defer log.Fatal("Stopping server")

	// Register our handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/host/", IdentityHandler)
	mux.HandleFunc("/workout/", api.GetWorkoutHandler)
	mux.HandleFunc("/", ui.IndexPage)

	// Register RESTfulHandlers
	mux.HandlerFunc("/metrics/bodyweight/", &api.Bodyweight{})

	// Start the server
	svrAddr := net.JoinHostPort(*addr, *port)
	if *cert != "" && *key != "" {
		http.ListenAndServerTLS(svrAddr, cert, key, mux)
	} else {
		http.ListenAndServe(svrAddr, mux)
	}
}
