package ui

import (
	"net/http"
	"text/template"
)

// Collection of initialized templates
var TemplateMap map[string]*html.Template

func init() {
	for name, path := range [][]string{
		{"index", "index.tmpl"},
	} {
		TemplateMap[name] = template.Must(template.New(name).Parse("html/" + path))
	}
}

// HomeHandler renders the homepage.
func HomeHandler(w http.ResponseWriter, r *http.Request) {

}
