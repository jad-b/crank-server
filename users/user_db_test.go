// +build test db

package users

import (
	"testing"

	"github.com/jad-b/torque"
)

var (
	username = "EzekielSparks"
	password = "BuildABetterButterRobot"
	user     = NewUserAccount(username, password)
)

func TestCreateUsersSchema(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	if err := CreateSchema(db); err != nil {
		t.Fatal(err)
	}
}

func TestCreateUserAuthTable(t *testing.T) {
	db := torque.Connect()
	defer db.Close()
	CreateSchema(db)

	if err := CreateTableUserAuth(db); err != nil {
		t.Fatal(err)
	}
}

func TestUserAuthCreate(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	if err := user.Create(db); err != nil {
		t.Fatal(err)
	}
	// Delete
	if err := user.Delete(db); err != nil {
		t.Fatal(err)
	}
}

func TestUserAuthRetrieve(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	if err := user.Create(db); err != nil {
		t.Fatal(err)
	}
	defer user.Delete(db)
	// Retrieve
	u := UserAuth{
		Username:  username,
		Superuser: true, // This bad value should get overridden
	}
	if err := u.Retrieve(db); err != nil {
		t.Fatal(err)
	}
	// Verify it's the right user.
	if u.Username != user.Username {
		t.Fatalf("Retrieved wrong user;\nGot: %s\nWanted: %s",
			u.Username, user.Username)
	}
	if u.Superuser {
		t.Fatalf("Failed to override Superuser being '%t'", u.Superuser)
	}
}

func TestUserAuthUpdate(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	if err := user.Create(db); err != nil {
		t.Fatal(err)
	}
	defer user.Delete(db)

	// Update
	u := *user // Make a copy of the shared struct
	// Give the user a new password
	u.PasswordHash, u.PasswordSalt, u.Cost = DefaultHash("new_password")
	// Disable the account - maybe we need the user to confirm
	u.Enabled = false
	u.Update(db)

	// Retrieve
	u2 := UserAuth{Username: username}
	if err := u2.Retrieve(db); err != nil {
		t.Fatal(err)
	}
	if u2.PasswordHash == user.PasswordHash {
		t.Fatalf("Password wasn't changed\n%s == \n%s\nExpected: %s",
			u2.PasswordHash, user.PasswordHash, u.PasswordHash)
	}
	if u2.PasswordHash != u.PasswordHash {
		t.Fatalf("Failed to update password\nGot: %s\nWanted: %s",
			u2.PasswordHash, u.PasswordHash)
	}
	if u2.Enabled {
		t.Fatal("Failed to disable user account")
	}
}

func TestUserAuthDelete(t *testing.T) {
	db := torque.Connect()
	defer db.Close()

	// Create
	if err := user.Create(db); err != nil {
		t.Fatal(err)
	}
	// Delete
	if err := user.Delete(db); err != nil {
		t.Fatal(err)
	}
	// Retrieve (and fail)
	if err := user.Retrieve(db); err == nil {
		t.Fatal("Failed to delete user")
	}
}
