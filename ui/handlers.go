/*Package ui is responsible for Crank's web UI.
 */
package ui

import (
	"html/template"
	"net/http"
	"path"
)

type bodyContent struct {
	Cowsay, Sponsor string
}

// HomePage renders a landing page.
func HomePage(w http.ResponseWriter, req *http.Request) {
	templatePath := path.Join("ui", "templates", "index.html")
	tmpl := template.Must(template.ParseFiles(templatePath))

	content := &bodyContent{Cowsay: "How's it going?", Sponsor: "Bill Nye: The Science Guy"}
	if err := tmpl.Execute(w, content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
