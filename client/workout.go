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
	u.Path = "/workout/"
	u.RawQuery = url.Values{
		// string => []string
		"timestamp": {timestamp.String()},
	}.Encode()

	res, err := http.Get(u.String())
	body, err := ioutil.ReadAll(res.Body)
	log.Printf("Response Body: %s", body)

	err = json.Unmarshal(body, w)
	if err != nil {
		log.Fatal("Failed to unmarshal workout: %s", err)
		log.Fatal("Response Body: %s", body)
	}
	log.Printf("Received this: %+v", w)
	return
}
