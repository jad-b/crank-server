// +build test db

package users

import (
	"log"
	"net/url"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	username = "EzekielSparks"
	password = "BuildABetterButterRobot"
	user     = NewUserAccount(username, password)
	testURL  = url.URL{
		Scheme: "http",
		Host:   "localhost:18000",
		// Generate a random user ID for working with
		Path: "users",
	}
	db *sqlx.DB
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestNewUserAccount(t *testing.T) {
	u := NewUserAccount(username, password)
	if u.PasswordHash == "" {
		t.Error("No password was created")
	}
	if u.PasswordHash == password {
		t.Error("Password wasn't hashed")
	}
	if u.PasswordSalt == "" {
		t.Error("No salt created for account")
	}
	if u.Cost == 0 {
		t.Error("No hashing-cost set on account")
	}
}
