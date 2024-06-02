package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"
)

func renderTemplate(w http.ResponseWriter, templateName string, data map[string]any) {
	//find all templates in templates folder
	templates, err := filepath.Glob("templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Read all templates into template set
	t, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//execute template
	t.ExecuteTemplate(w, templateName, data)
}
