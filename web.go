package torque

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// Common HTTP constants
const (
	HeaderAuthorization = "Authorization"
	HeaderAuthenticate  = "WWW-Authenticate"

	MimeJSON = "application/json"

	// Scheme dictates http vs. https (or anything else, I suppose...)
	// TODO switch to https
	Scheme = "http"
)

var (
	// ValidTimestamps are all approved datetime formats in Torque
	// See RFC 1123
	ValidTimestamps = []string{
		time.RFC822,
		time.RFC850,
		time.ANSIC,
	}
)

// A RESTfulResource knows how to represent itself in a RESTful API.
type RESTfulResource interface {
	GetResourceName() string
}

// The RESTfulHandler interface models a resource that supports RESTful
// interactions. It supports GET, POST, PUT, and DELETE operations, as well as
// the ServeHTTP method required for HTTP handlers.
// TODO(jdb) Add an `error` return to each Handle* method, along with custom
// error types. Let the RequestHandler type generically deal with different
// errors.
type RESTfulHandler interface {
	RESTfulResource
	HandleGet(http.ResponseWriter, *http.Request)
	HandlePost(http.ResponseWriter, *http.Request)
	HandlePut(http.ResponseWriter, *http.Request)
	HandleDelete(http.ResponseWriter, *http.Request)
}

// RequestHandler wraps HTTP request handling in Torque's best practices for
// logging and error handling. Or will - right now I'm just playing around.
// TODO(jdb) Add an `error` return value
type RequestHandler func(http.ResponseWriter, *http.Request)

// TODO(jdb) Handle returned errors by type
func (fn RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log incoming requests
	LogRequest(r)
	fn(w, r)
}

// RouteRequest returns the corresponding method based on the incoming
// request's HTTP method.
//
// Example Usage:
//   http.HandleFunc("/foo", RouteRequest(metrics.Bodyweight{}))
func RouteRequest(rr RESTfulHandler) func(http.ResponseWriter, *http.Request) {
	// Accept RESTfulHandler, return f(w, req)
	return func(w http.ResponseWriter, req *http.Request) {
		switchByMethod(rr, req)
	}
}

func switchByMethod(rr RESTfulHandler, req *http.Request) RequestHandler {
	switch req.Method {
	case "GET":
		return rr.HandleGet
	case "POST":
		return rr.HandlePost
	case "PUT":
		return rr.HandlePut
	case "DELETE":
		return rr.HandleDelete
	default:
		return func(writer http.ResponseWriter, request *http.Request) {
			http.Error(
				writer,
				fmt.Sprintf("%s is not a support HTTP method for this resource",
					request.Method),
				http.StatusMethodNotAllowed)
		}
	}
}

// LogRequest idempotently writes the http.Request to the default Logger.
func LogRequest(r *http.Request) {
	b, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(string(b))
}

// LogRequestThenError dumps the request into log output and returns an error. It is really
// only good as a placeholder, which is why it returns an 501 Not Implemented error.
func LogRequestThenError(w http.ResponseWriter, r *http.Request) {
	LogRequest(r)
	http.Error(w,
		"Your request was logged, but no functionality exists at this endpoint.",
		http.StatusNotImplemented)
}

// LogResponse idempotently writes the http.Response to the default Logger.
func LogResponse(resp *http.Response) {
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(b))
}

// GetOrCreateTimestamp ensures a timestamp is attached to the Request. First it looks for
// a Query field "timestamp". Failing that, it returns the current time.
// Query.
func GetOrCreateTimestamp(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	// Attempt to parse
	if &queryTime == nil {
		return time.Now(), nil
	}
	return ParseTimestamp(queryTime)
}

// ParseTimestamp applies all valid timestamps to the string value.
func ParseTimestamp(value string) (time.Time, error) {
	var err error
	for _, timeFmt := range ValidTimestamps {
		// Try to parse
		t, err := time.Parse(timeFmt, value)
		if err == nil { // If successful, return
			return t, nil
		}
	}
	return time.Time{}, err
}

// ReadJSONRequest reads the body of a request directly into a given struct.
func ReadJSONRequest(req *http.Request, v interface{}) error {
	return json.NewDecoder(req.Body).Decode(v)
}

// ReadJSONResponse unmarshals the http.Response.Body into a struct.
func ReadJSONResponse(resp *http.Response, v interface{}) error {
	return json.NewDecoder(resp.Body).Decode(v)
}

// WriteJSON writes the value v to the http response stream as json with standard
// json encoding.
// Stolen from github.com/docker/docker/api/server
func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("Failed to encode as json:\n\t%v\nSending %d", v,
			http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// WriteOkayJSON encodes v to the HTTP response with a 200 OK status code.
func WriteOkayJSON(w http.ResponseWriter, v interface{}) {
	WriteJSON(w, http.StatusOK, v)
}

// HTTPError wraps http.Error but handles marshalling the error into a JSON
// string
func HTTPError(w http.ResponseWriter, e error, code int) {
	// Marhsall struct into a JSON string
	errorJSON, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		log.Printf("Trouble marshalling this error: %s.\nThe user will receive a generic %s", e.Error(), genericErrorJSON)
		errorJSON = genericErrorJSON
	}
	http.Error(w, string(errorJSON), code)
}
