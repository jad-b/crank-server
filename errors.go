package torque

import "encoding/json"

type genericError struct {
	msg string `json:"error"`
}

var genericErrorJSON, _ = json.MarshalIndent(genericError{"Something went wrong."}, "", "\t")

// ErrorResponse is intended for HTTP response bodies.
type ErrorResponse struct {
	Message string `json:"error"`
}

// Error returns the internal Error message.
func (e ErrorResponse) Error() string {
	return e.Message
}
