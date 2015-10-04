package client

import "testing"

type UserIDStruct struct {
	UserID int
}

func TestSetUserID(t *testing.T) {
	s := UserIDStruct{}
	SetUserID(&s, 17)
	if s.UserID != 17 {
		t.Error("Failed to set UserID in struct")
	}
}

func TestSetUserIDNoOverwrite(t *testing.T) {
	s := UserIDStruct{42}
	SetUserID(&s, 17)
	if s.UserID != 42 {
		t.Error("Overwrote original value")
	}

}
