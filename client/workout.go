package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jad-b/crank/crank"
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
		"timestamp": {timestamp.Format(time.RFC3339)},
	}.Encode()
	log.Printf("Client: Requesting workout at time %s", timestamp.Format(time.RFC3339))

	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read response body:\n%s", res.Body)
	}
	log.Printf("Client: Response Body: %s", body)

	// ??? We need to pass a reference to the pointer?
	err = json.Unmarshal(body, &w)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Client: Received this: %+v", w)
	return
}
