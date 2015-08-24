package torque

import (
	"fmt"
	"net/http"
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
