package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jad-b/torque/ui"
)

var (
	addr = flag.String("addr", "127.0.0.1:18001", "Host:port of Torque UI server")
)

func main() {
	// Serve static assets
	http.HandleFunc("/", ui.ReloadHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets"))))

	log.Print("Serving on ", *addr)
	http.ListenAndServe(*addr, nil)
	log.Fatal("Stopping server")
}
