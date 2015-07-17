package api

import (
	"fmt"

	"github.com/jad-b/torque"
)

// The RESTfulHandler interface models a resource that supports RESTful
// interactions. It supports GET, POST, PUT, and DELETE operations, as well as
// the ServeHTTP method required for HTTP handlers.
type RESTfulHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	Put(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	http.Handler
}

// The RESTfulRouter supports routing incoming HTTP requests to the correct
// handler by way of RouteMethod. It is intended for embedding in structs which
// require RESTful method routing.
type RESTfulRouter struct{}

// Delete is a no-op
func (rr *RESTfulRouter) Delete(w http.ResponseWriter, req *http.Request) {
	web.LogRequestThenError(w, req)
}

// Get is a no-op
func (rr *RESTfulRouter) Get(w http.ResponseWriter, req *http.Request) {
	web.LogRequestThenError(w, req)
}

// Post is a no-op
func (rr *RESTfulRouter) Post(w http.ResponseWriter, req *http.Request) {
	web.LogRequestThenError(w, req)
}

// Put is a no-op
func (rr *RESTfulRouter) Put(w http.ResponseWriter, req *http.Request) {
	web.LogRequestThenError(w, req)
}

// RouteMethod returns the corresponding method based on the incoming
// request's HTTP method.
func (rr *RestfulRouter) RouteMethod(req *http.Request) func(http.ResponseWriter, *http.Request) {
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

func (rr *RestfulRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	return rr.RouteMethod(req)
}
