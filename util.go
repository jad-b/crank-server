package torque

import (
	"encoding/json"
	"log"
	"strings"
)

// SlashJoin performs a strings.Join using '/' as a separator.
func SlashJoin(args ...string) string {
	return strings.Join(args, "/")
}

// PrettyJSON pretty-prints JSON. If an error occurs, you'll get back an empty,
// but valid, JSON structure.
func PrettyJSON(v interface{}) string {
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Print(err)
		return "{}"
	}
	return string(s)
}
