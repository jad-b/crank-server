package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	addr = flag.String("addr", "127.0.0.1:18001", "Host:port of Torque UI server")
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", HomeHandler)

	http.ListenAndServe(*addr, mux)
	log.Fatal("Stopping server")
}
