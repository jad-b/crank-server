package main

import (
	"encoding/json"
	"fmt"
	"github.com/jad-b/crank/crank"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Print("Starting server")
	http.HandleFunc("/host/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("host is %s", req.Host)
		fmt.Fprintf(w, "%s, this is me.", req.Host)
	})
	http.HandleFunc("/workout/", func(w http.ResponseWriter, req *http.Request) {
		log.Print("Requested /workout/")
		// Write a stub workout to the reesponse
		now := time.Now()
		json.NewEncoder(w).Encode(&crank.Workout{
			Timestamp: now,
			Comment:   fmt.Sprintf("Time is %s", now.String()),
		})
	})
	http.ListenAndServe(":8000", nil)
	defer log.Fatal("Stopping server")
}
