package users

import "time"

// UserPIISQL is SQL for creating the auth.UserPII table.
const UserPIISQL = `
CREATE TABLE auth.UserPII (
  "id" integer PRIMARY KEY,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "email" text NOT NULL UNIQUE,
  "last_modified" timestamp(0) with time zone NOT NULL
);`

// UserPII contains sensitive user information which could identify them in the
// real-world. Most of the time, you should only be accessing this if you need
// to display the user's profile or send them an email. Or swag.
//
// PII (Personally Identifiable Information) has been extracted from the rest
// of the user tables. This is in the hope of preventing anything stupid from
// accidentally happening to it (log dumps, over-eager SELECT statements).
type UserPII struct {
	// Primary key. Used for linking User tables.
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	LastModified time.Time `json:"last_modified"`
}
