package ui

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/jad-b/torque"
)

var (
	// TemplateReqs	maps page names to template dependecies
	templateReqs = map[string][]string{
		"main":  []string{"main", "head", "header", "sidebar", "footer", "nav", "scripts", "base"},
		"index": []string{"index", "head", "header", "sidebar", "footer", "nav", "scripts", "base"},
	}
	// TemplateMap is a collection of initialized templates
	TemplateMap = make(map[string]*template.Template)
)

// WktHandler converts Workouts in JSON into the .wkt format
func WktHandler(w http.ResponseWriter, r *http.Request) {
	var workout Workout // Unmarshal workout from request
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if s, err := torque.PrettyJSON(workout); err == nil {
		log.Print(s) // Log for debugging
		wkt, err := WorkoutToWkt(&workout)
		if err == nil {
			log.Print(wkt)
			err = json.NewEncoder(w).Encode(struct{ Wkt string }{Wkt: wkt})
			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

// ReloadHandler always loads the template on every response
func ReloadHandler(templateDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine which template to render
		path := r.URL.Path[1:] // Drop leading '/'
		var name string
		if len(path) == 0 {
			name = "index"
		} else {
			name = path
		}
		log.Print("Serving ", name)

		// Load template
		reqs := prepFilepaths(templateReqs[name], templateDir)
		t := template.Must(template.ParseFiles(reqs...))
		// Print out all associated templates //log.Printf("Rendering template %s", t.Name())
		//for _, tPtr := range t.Templates() {
		//log.Printf("\t%s", tPtr.Name())
		//}

		// Render template, using the specified base
		if err := t.ExecuteTemplate(w, "base.html", nil); err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})
}

// LoadTemplates parses all the templates into memory
func LoadTemplates(tmplDir string) {
	log.Printf("Loading html from %s", tmplDir)
	for name, files := range templateReqs {
		tmpls := prepFilepaths(files, tmplDir)
		TemplateMap[name] = template.Must(template.ParseFiles(tmpls...))
	}
}

func prepFilepaths(tmplReqs []string, tmplDir string) (files []string) {
	// Lookup template deps
	tmpls := make([]string, len(tmplReqs))
	// Prefix asset dir path to every file; suffix with '.tmpl'
	for i := range tmplReqs {
		tmpls[i] = torque.SlashJoin(tmplDir, tmplReqs[i]) + ".html"
	}
	return tmpls
}
