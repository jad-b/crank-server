package torque

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/jmoiron/sqlx"
)

// A collection of test utilities.

// Connect sets up a DB connection for the test.
// Only tests should need to frequently setup DB connections.
// This also makes the assumption you want to use the '-psql-conf' flag.
//
// Don't forgot to `defer db.Close()`
func Connect() *sqlx.DB {
	// Setup our database connection
	pgConf := LoadPostgresConfig(*PsqlConf)
	return OpenDBConnection(pgConf)
}

// DieOnError fatally logs an error, if it's real.
func DieOnError(t *testing.T, err error) {
	if err != nil {
		_, f, l, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: %s", f, l, err.Error())
	}
}

// DumpRecordedResponse writes a http.ResponseRecorder to a string.
// It takes care of resetting the Body for future ops.
func DumpRecordedResponse(r *httptest.ResponseRecorder) string {
	bodyCopy := *r.Body // Copy the body
	var headerBuffer bytes.Buffer
	r.HeaderMap.Write(&headerBuffer)
	return fmt.Sprintf("HTTP Recorded Response (%d):\n%s\n%s\n",
		r.Code, headerBuffer.String(), bodyCopy.String())
}
