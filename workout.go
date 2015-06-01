package crank

import (
	"log"
	"net/http"
	"net/url"
)

// GetWorkout retrieves a workout from the server.
func GetWorkout(base String, timestamp time.Time) (w *Workout, err error) {
	url, err := url.Parse(base)
	url.Path = "/workout/"
	&url.URL{Scheme: "http", Host: serverUrl, Path: "/workout/"}
	url.RawQuery = &url.Values{
		// string => []string
		"timestamp": {timestamp},
	}.Encode()
	res, err := http.Get(url.String())
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, w)
	log.Printf("Received this: %+v", w)
	return
}
