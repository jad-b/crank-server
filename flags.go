package torque

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
)

// CommandLineActor is capable of parsing and acting upon commmand-line arguments
type CommandLineActor interface {
	DBActor
	ParseFlags(action string, args []string)
}

// ActOnDB requests the actor perform it's correct method against the database.
func ActOnDB(actor DBActor, action string, db *sqlx.DB) error {
	switch action {
	case "create":
		return actor.Create(db)
	case "retrieve":
		return actor.Retrieve(db)
	case "update":
		return actor.Update(db)
	case "delete":
		return actor.Delete(db)
	default:
		return fmt.Errorf("%s is an invalid action", action)
	}
}

// ActOnWebServer requests the actor perform it's correct method against a web
// server.
func ActOnWebServer(action, serverURL string) (*http.Response, error) {
	switch action {
	case "create":
		return nil, nil
	case "retrieve":
		return nil, nil
	case "update":
		return nil, nil
	case "delete":
		return nil, nil
	default:
		return nil, fmt.Errorf("%s is an invalid action", action)
	}
}

// TimestampFlag is a custom command-line flag for accepting timestamps
type TimestampFlag time.Time

func (ts *TimestampFlag) String() string {
	return time.Time(*ts).String()
}

// Set reads the raw string value into a TimestampFlag, or dies
// trying...actually it just returns nil.
func (ts *TimestampFlag) Set(value string) error {
	t, err := ParseTimestamp(value)
	if err != nil {
		return err
	}
	*ts = TimestampFlag(t)
	return nil
}

// HostPortFlag handles converting a "host:port" string into a full-fledged
// url.URL
type HostPortFlag url.URL

// String calls url.URL's String() method
func (hpf *HostPortFlag) String() string {
	u := url.URL(*hpf)
	return u.String()
}

// Set parses the host:port string into a valid url.URL
func (hpf *HostPortFlag) Set(value string) error {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return err
	}
	hp := net.JoinHostPort(host, port)
	*hpf = HostPortFlag(url.URL{
		Host: hp,
	})
	return nil
}
