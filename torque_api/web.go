package torque

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Common HTTP constants
const (
	HeaderAuthorization = "Authorization"
	HeaderAuthenticate  = "WWW-Authenticate"
	HeaderContentType   = "Content-Type"

	MimeJSON = "application/json"

	// Scheme dictates http vs. https (or anything else, I suppose...)
	// TODO switch to https
	Scheme = "http"
)

var (
	// ValidTimestamps are all approved datetime formats in Torque
	// See RFC 1123
	ValidTimestamps = []string{
		time.RFC3339Nano,
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
// error types. Let the SmartHandler type generically deal with different
// errors.
type RESTfulHandler interface {
	RESTfulResource
	HandleGet(http.ResponseWriter, *http.Request)
	HandlePost(http.ResponseWriter, *http.Request)
	HandlePut(http.ResponseWriter, *http.Request)
	HandleDelete(http.ResponseWriter, *http.Request)
}

// ErrorableHandler extends the http.Handler signature to return an error.
type ErrorableHandler func(w http.ResponseWriter, r *http.Request) error

// ErrorHandler wraps ErrorableHandler's in error-handling code.
type ErrorHandler struct {
	context *Context
	fn      ErrorableHandler
}

func (eh *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := eh.fn(w, r); err != nil {
		switch e := err.(type) {
		case ErrorResponse:
			http.Error(w, e.Message, e.StatusCode)
		default:
			http.Error(w, e.Error(), http.StatusBadRequest)
		}
	}
}

// SmartHandler provides Torque best practices on top of HTTP request handling.
type SmartHandler func(http.ResponseWriter, *http.Request)

// ServeHTTP applies middleware steps to requests.
func (sh SmartHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Print("Handling HTTP request smartly...")
	LogRequest(req)
	sh(w, req)
}

// RouteRequest returns the corresponding method based on the incoming
// request's HTTP method.
//
// Example Usage:
//   http.HandleFunc("/foo", RouteRequest(metrics.Bodyweight{}))
func RouteRequest(rr RESTfulHandler) SmartHandler {
	// Closure associates the RESTfulHandler with the normal HandleFunc
	// function signature
	fn := func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			rr.HandleGet(w, req)
		case "POST":
			rr.HandlePost(w, req)
		case "PUT":
			rr.HandlePut(w, req)
		case "DELETE":
			rr.HandleDelete(w, req)
		default:
			func(writer http.ResponseWriter, request *http.Request) {
				http.Error(
					writer,
					fmt.Sprintf("%s is not a support HTTP method for this resource",
						request.Method),
					http.StatusMethodNotAllowed)
			}(w, req)
		}
	}
	// Wrap all RESTful Handler routing in Torque's request handling
	return SmartHandler(fn)
}

// LogRequest idempotently writes the http.Request to the default Logger.
func LogRequest(req *http.Request) {
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Print("Request: ", err)
		// Try to write the whole thing ourselves
		var b bytes.Buffer
		err = req.Write(&b)
		log.Print(b.String())
		if err != nil {
			log.Print("Failed to write Request; ", err)
		}
		return
	}
	log.Print("Request: ", string(b))
}

// ServeDir returns a handler for serving a file directory
func ServeDir(prefix, dirpath string) http.Handler {
	fs := http.FileServer(http.Dir(dirpath))
	if prefix != "" {
		return http.StripPrefix(prefix, fs)
	}
	return fs
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

// GetTimestampQuery ensures a timestamp is attached to the Request. First it looks for
// a Query field "timestamp". Failing that, it returns the current time.
// Query.
func GetTimestampQuery(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	return ParseTimestamp(queryTime)
}

// SetTimestampQuery attaches a timestamp query parameter to the request.
func SetTimestampQuery(u *url.URL, t time.Time) {
	stampString := Stamp(t)
	q := u.Query()
	q.Set("timestamp", stampString)
	u.RawQuery = q.Encode()
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

// WriteOkayJSON encodes v to the HTTP response with a 200 OK status code.
func WriteOkayJSON(w http.ResponseWriter, v interface{}) error {
	return WriteJSON(w, http.StatusOK, v)
}

// WriteJSON writes the value v to the http response stream as json with standard
// json encoding.
// Stolen from github.com/docker/docker/api/server
func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	if v != nil {
		jason, err := PrettyJSON(v)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to encode as json:\n\t%v")
			http.Error(w, errMsg, http.StatusInternalServerError)
			return err
		}
		w.Header().Set(HeaderContentType, MimeJSON)
		io.WriteString(w, jason)
		return nil
	}
	return errors.New("No interface provided")
}

// HTTPError wraps http.Error but handles marshalling the error into a JSON
// string
func HTTPError(w http.ResponseWriter, e error, code int) {
	// Marhsall struct into a JSON string
	var errResp ErrorResponse
	var ok bool
	if errResp, ok = e.(ErrorResponse); !ok { // Convert only if necessary
		errResp = ErrorResponse{Message: e.Error()}
	}
	errorJSON, err := json.MarshalIndent(errResp, "", "\t")
	if err != nil {
		log.Printf("Trouble marshalling this error: %s.\nThe user will receive a generic %s", e.Error(), genericErrorJSON)
		errorJSON = genericErrorJSON
	}
	http.Error(w, string(errorJSON), code)
}
