package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/ui"
)

var (
	addr = flag.String("addr", "127.0.0.1:18001", "Host:port of Torque UI server")
	// Assets assigns the location of static assets
	Assets = flag.String("assets", "./ui/assets", "Location of html/css/js assets")
)

func main() {
	mux := http.NewServeMux()
	templateDir := *Assets + "/html"

	// Serve static assets
	mux.Handle("/assets/", torque.ServeDir("/assets/", "ui/assets"))
	// Chrome will *not* let the favicon be anywhere but '/'. Oh well.
	mux.Handle("/favicon.ico", torque.ServeDir("", "ui/assets/img/"))
	mux.Handle("/spec/", torque.ServeDir("/spec/", "."))

	// Handlers that do something
	mux.Handle("/wkt", LogRequest(http.HandlerFunc(ui.WktHandler)))
	mux.Handle("/", LogRequest(ui.ReloadHandler(templateDir)))

	log.Print("Serving on ", *addr)
	http.ListenAndServe(*addr, mux)
	log.Fatal("Stopping server")
}

// LogRequest outputs the incoming request for debugging purposes
func LogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b []byte
		var err error
		if b, err = httputil.DumpRequest(r, true); err != nil {
			log.Print(err)
		}
		log.Print(string(b))

		h.ServeHTTP(w, r)
	})
}
