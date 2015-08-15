// The users.session package exposes session management utilities, (soon to be)
// including auth.
package users

import (
	"time"
)

type UserSession struct {
	SessionKey string    `json:session_key`
	UserID     int       `json:user_id`
	LoginTime  time.Time `json:login_time`
	LastSeen   time.Time `json:last_seen`
}
