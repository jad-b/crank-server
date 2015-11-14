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
	// TemplateReqs	maps page names to template dependecies
	templateReqs = map[string][]string{
		"index": []string{"index", "footer", "base"},
	}
	// TemplateMap is a collection of initialized templates
	TemplateMap = make(map[string]*template.Template)
	// Assets assigns the location of static assets
	Assets      = flag.String("assets", "./ui/assets", "Location of html/css/js assets")
	templateDir = torque.SlashJoin(*Assets, "html")
)

func init() {
	log.Printf("Loading html from %s", templateDir)
	for name, files := range templateReqs {
		tmpls := prepFilepaths(files, templateDir)
		TemplateMap[name] = template.Must(template.ParseFiles(tmpls...))
	}
}

func prepFilepaths(tmplReqs []string, tmplDir string) (files []string) {
	// Lookup template deps
	tmpls := make([]string, len(tmplReqs))
	// Prefix asset dir path to every file; suffix with '.tmpl'
	for i := range tmplReqs {
		tmpls[i] = torque.SlashJoin(tmplDir, tmplReqs[i]) + ".tmpl"
	}
	return tmpls
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
	// Determine which template to render
	path := r.URL.Path[1:] // Drop leading '/'
	var name string
	if len(path) == 0 {
		name = "index"
	} else {
		name = path
	}

	// Log request for debugging purposes
	b, _ := httputil.DumpRequest(r, false)
	log.Print(string(b))

	// Load template
	reqs := prepFilepaths(templateReqs[name], templateDir)
	t := template.Must(template.ParseFiles(reqs...))
	// Print out all associated templates
	log.Printf("Rendering template %s", t.Name())
	for _, tPtr := range t.Templates() {
		log.Printf("\t%s", tPtr.Name())
	}

	// Render template, using the specified base
	err := t.ExecuteTemplate(w, "base.tmpl", nil)
	if err != nil {
		log.Print(err)
	}
}
