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
	if user.PasswordHash == "" {
		t.Error("No password was created")
	}
	if user.PasswordHash == password {
		t.Error("Password wasn't hashed")
	}
	if user.PasswordSalt == "" {
		t.Error("No salt created for account")
	}
	if user.Cost == 0 {
		t.Error("No hashing-cost set on account")
	}
}

func TestPasswordValidation(t *testing.T) {
	if !user.ValidatePassword(password) {
		t.Fatal("Password wasn't accepted")
	}
}
