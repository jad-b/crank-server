// +build test metrics

package redteam

import (
	"net/url"
	"testing"
	"time"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/client"
	"github.com/jad-b/torque/metrics"
)

func TestPostingBodyweight(t *testing.T) {
	// Connect to the server
	tAPI := client.NewTorqueAPI(torqueAddr.String())

	// Authenticate our user
	err := tAPI.Authenticate(username, password)
	torque.DieOnError(t, err)

	// Create BW record
	now := time.Now()
	bw := metrics.Bodyweight{
		UserID:    tAPI.User.ID,
		Timestamp: now,
		Weight:    181.2,
		Comment:   "I made this up",
	}

	// Post to server
	_, err = tAPI.Post(&bw)
	torque.DieOnError(t, err)
	t.Log("POST'd Bodyweight record")

	// Retrieve record
	resp, err := tAPI.Get(&bw, url.Values{"timestamp": []string{torque.Stamp(now.Truncate(time.Second))}})
	torque.DieOnError(t, err)
	t.Log("GET'd Bodyweight record")

	if resp.StatusCode != 200 {
		torque.LogResponse(resp)
		t.Fatal("Non-200 status code returned")
	}

	// Read bodyweight record from response
	var bw2 metrics.Bodyweight
	err = torque.ReadJSONResponse(resp, &bw2)
	torque.DieOnError(t, err)
	if bw2.Weight != bw.Weight || bw2.Comment != bw2.Comment {
		t.Fatalf("Not equal:\n%#v\n%#v", bw, bw2)
	}
	t.Log("POST & GET successful")
}
