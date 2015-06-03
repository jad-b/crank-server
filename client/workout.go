package client

import (
	"encoding/json"
	"github.com/jad-b/crank/crank"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// GetWorkout retrieves a workout from the server.
func GetWorkout(base string, timestamp time.Time) (w *crank.Workout, err error) {
	u, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}
	u.Path = "/workout/"
	u.RawQuery = url.Values{
		// string => []string
		"timestamp": {timestamp.String()},
	}.Encode()
	log.Printf("Requesting workout at time %s", timestamp)

	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read response body:\n%s", res.Body)
	}
	log.Printf("Response Body: %s", body)

	// ??? We need to pass a reference to the pointer?
	err = json.Unmarshal(body, &w)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Received this: %+v", w)
	return
}
