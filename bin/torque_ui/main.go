package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jad-b/torque/ui"
)

var (
	addr = flag.String("addr", "127.0.0.1:18001", "Host:port of Torque UI server")
)

func main() {
	mux := mux.NewRouter()
	// Serve static assets
	mux.HandleFunc("/", ui.ReloadHandler)
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets"))))

	log.Print("Serving on ", *addr)
	http.ListenAndServe(*addr, mux)
	log.Fatal("Stopping server")
}
