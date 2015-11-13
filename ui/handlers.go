package ui

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
	"github.com/jad-b/torque"
)

var (
	// TemplateMap is a collection of initialized templates
	TemplateMap = make(map[string]*template.Template)
	// Assets assigns the location of static assets
	Assets      = flag.String("assets", "./ui/assets", "Location of html/css/js assets")
	templateDir = torque.SlashJoin(*Assets, "html")
)

func init() {
	log.Printf("Loading html from %s", templateDir)
	for name, files := range map[string][]string{
		"index": []string{"index.tmpl"},
	} {
		// Prefix 'html/' to every file
		tmpls := make([]string, len(files))
		for i := range files {
			tmpls[i] = torque.SlashJoin(templateDir, files[i])
		}
		TemplateMap[name] = template.Must(template.ParseFiles(tmpls...))
	}
}

// Routes returns the routes available for the UI
func Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	return r
}

// HomeHandler renders the homepage.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	TemplateMap["index"].Execute(w, nil)
}

// ReloadHandler always loads the template on every response
func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := httputil.DumpRequest(r, true)
	log.Print(string(b))
	// Load template
	t := template.Must(template.ParseFiles(torque.SlashJoin(templateDir, "dev.tmpl")))
	log.Print("Rendering template")
	t.Execute(w, nil)
}
