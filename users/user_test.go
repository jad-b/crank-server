package users

import (
	"net/http/httptest"
	"testing"
)

func TestAuthentication(t *testing.T) {
	serverURL, username, password := "https://localhost", "JohnFritz", "gazebo"
	// Create a request for authentication
	req, err := buildAuthenticationRequest(serverURL, username, password)
	w := httptest.NewRecorder()

	// See how we handle the authentication
	HandleAuthentication(w, req)

}
