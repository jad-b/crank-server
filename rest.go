package torque

import (
	"fmt"
	"net/http"
)

// The RESTfulHandler interface models a resource that supports RESTful
// interactions. It supports GET, POST, PUT, and DELETE operations, as well as
// the ServeHTTP method required for HTTP handlers.
type RESTfulHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	Put(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

// RouteRequest returns the corresponding method based on the incoming
// request's HTTP method.
//
// Example Usage:
//   rr := &Bodyweight{}
//   http.HandleFunc("/foo", RouteMethod(rr))
func RouteRequest() func(http.ResponseWriter, *http.Request) {
	return func(http.ResponseWriter, *http.Request) {
		switch req.Method {
		case "GET":
			return rr.Get
		case "POST":
			return rr.Post
		case "PUT":
			return rr.Put
		case "DELETE":
			return rr.Delete
		default:
			return func(writer http.ReponseWriter, request *http.Request) {
				return http.Error(
					writer,
					fmt.Sprintf("%s is not a support HTTP method for this resource", req.Method),
					http.StatusMethodNotAllowed)
			}
		}
	}
}
