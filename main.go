package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting server")
	http.HandleFunc("/host/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("host is %s", req.Host)
		fmt.Fprintf(w, "%s, this is me.", req.Host)
	})
	http.ListenAndServe(":8000", nil)
	defer log.Fatal("Stopping server")
}
