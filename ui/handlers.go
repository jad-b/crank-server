/*
Package ui is responsible for Crank's web UI.

A couple naming conventions are observed. Files which can be rendered stand-lone as
valid HTML end with '.html', while files which contain only template definitions
end with '.tmpl'. This helps delineate between parent and child template files.

Additionally, handlers for UI web pages are named "<Something>Page". This distinguishes
them from API endpoints.
*/
package ui

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

const (
	templateDir string = "ui/templates"
)

var (
	basePath      = path.Join(templateDir, "base.html")
	baseTemplate  = template.Must(template.New(basePath))
	indexTemplate = baseTemplate.Clone().Parse(path.Join(templateDir, "index.tmpl"))
)

// RenderPage renders
func RenderPage(w http.ResponseWriter, req *http.Request, data interface{},
	filenames ...string) {
	// Prepend our base template
	templatePaths := append([]string{basePath}, filenames...)
	// Load and compile our templates into a single Template object
	tmpl, _ := template.ParseFiles(templatePaths...)
}

// IndexPage renders a landing page. It's pretty stupid.
func IndexPage(w http.ResponseWriter, req *http.Request) {
	index := path.Join(templateDir, "index.tmpl")
	data := struct {
		Cowsay, Sponsor string
	}{
		Cowsay:  "How's it going?",
		Sponsor: "Bill Nye: The Science Guy",
	}
	// An error here means we had trouble actually rendering the template(s)
	if err := indexTemplate.Execute(w, data); err != nil {
		// web.LogHTTPError(w, err)
		log.Print(err.Error())
	}
}
