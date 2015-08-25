package torque

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// A RESTfulResource knows how to represent itself in a RESTful API.
type RESTfulResource interface {
	GetResourceName() string
}

// The RESTfulHandler interface models a resource that supports RESTful
// interactions. It supports GET, POST, PUT, and DELETE operations, as well as
// the ServeHTTP method required for HTTP handlers.
type RESTfulHandler interface {
	RESTfulResource
	HandleGet(http.ResponseWriter, *http.Request)
	HandlePost(http.ResponseWriter, *http.Request)
	HandlePut(http.ResponseWriter, *http.Request)
	HandleDelete(http.ResponseWriter, *http.Request)
}

// RouteRequest returns the corresponding method based on the incoming
// request's HTTP method.
//
// Example Usage:
//   rr := &Bodyweight{}
//   http.HandleFunc("/foo", RouteMethod(rr))
func RouteRequest(rr RESTfulHandler) func(http.ResponseWriter, *http.Request) {
	// Accept RESTfulHandler, return f(w, req)
	return func(w http.ResponseWriter, req *http.Request) {
		switchByMethod(rr, req)
	}
}

func switchByMethod(rr RESTfulHandler, req *http.Request) func(http.ResponseWriter, *http.Request) {
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

// NopCloser satisfies the ReadCloser interface
type NopCloser struct{ io.Reader }

// Close does nothing but satisfy the Closer interface
func (nc NopCloser) Close() error { return nil }

// RequestToBuffer writes the HTTP request to a buffer for printing
func RequestToBuffer(req *http.Request) *bytes.Buffer {
	var buf *bytes.Buffer
	req.Write(buf)
	// Reset Request Body, since it's an io.ReadCloser that's been read
	req.Body = NopCloser{buf}
	return buf
}

// LogRequestThenError dumps the request into log output and returns an error. It is really
// only good as a placeholder, which is why it returns an 501 Not Implemented error.
func LogRequestThenError(w http.ResponseWriter, req *http.Request) {
	log.Printf("Incoming request:\n%s", RequestToBuffer(req).String())
	http.Error(w,
		"Your request was logged, but no functionality exists at this endpoint.",
		http.StatusNotImplemented)
}

// ReadBody extracts the body from the HTTP request. If there is an error, it
// writes it back to the response.
func ReadBody(w http.ResponseWriter, req *http.Request) (b []byte) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return nil
	}
	return b
}

// ReadBodyTo reads the body of a request directly into a given struct.
func ReadBodyTo(w http.ResponseWriter, req *http.Request, v interface{}) error {
	return json.NewDecoder(req.Body).Decode(v)
}

// Stamp ensures a timestamp is attached to the Request. First it looks for
// a Query field "timestamp". Failing that, it returns the current time.
// Query.
func Stamp(req *http.Request) (t time.Time, err error) {
	queryTime := req.URL.Query().Get("timestamp")
	if &queryTime == nil {
		return time.Now(), nil
	}
	return time.Parse(time.RFC3339, queryTime)
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
