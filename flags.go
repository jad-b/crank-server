package torque

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// CommandLineActor is capable of parsing and acting upon commmand-line arguments
type CommandLineActor interface {
	DBActor
	RESTfulClient
	ParseFlags(action string, args []string)
}

// ActOnDB requests the actor perform it's correct method against the database.
func ActOnDB(actor DBActor, action string, conn *sql.DB) error {
	switch action {
	case "create":
		return actor.Create(conn)
	case "retrieve":
		return actor.Retrieve(conn)
	case "update":
		return actor.Update(conn)
	case "delete":
		return actor.Delete(conn)
	default:
		return fmt.Errorf("%s is an invalid action", action)
	}
}

// ActOnWebServer requests the actor perform it's correct method against a web
// server.
func ActOnWebServer(actor RESTfulClient, action, serverURL string) (*http.Response, error) {
	switch action {
	case "create":
		return actor.HTTPPost(serverURL)
	case "retrieve":
		return actor.HTTPGet(serverURL)
	case "update":
		return actor.HTTPPut(serverURL)
	case "delete":
		return actor.HTTPDelete(serverURL)
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
	// TODO Change to a  list of valid timestamp formats
	MyTimeFormat := "2006Jan21504"
	t, err := time.Parse(MyTimeFormat, value)
	if err != nil {
		return err
	}
	*ts = TimestampFlag(t)
	return nil
}
