// +build test api,db

package users

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUserCreation(t *testing.T) {
	username, Password := "newuser", "newuserpassword"
	user := &UserAuth{}

	// Create request
	u := url.URL{Scheme: torque.Scheme, Host: "localhost", Path: user.GetResourceName()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	user.HandlePost(resp, req)

	if resp.Code != 200 {
		t.Fatal(resp.Body.String())
	}
}
