package redteam

import (
	"github.com/jad-b/torque"
	"github.com/jad-b/torque/metrics"
	"github.com/jad-b/torque/users"
)

/*
metrics provides interactions with the metrics/ API
*/

// CreateBodyweightRecord POSTs a new bodyweight record to the server.
// Assume the user has been authenticated.
func CreateBodyweightRecord(api *torque.API, u *users.UserAuth) {
	bw := metrics.GenerateBodyweight(180.0, "If you're not growing, you're dying")
	respBW, err := api.Post(user, bw)

}
