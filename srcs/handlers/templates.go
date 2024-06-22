package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"
)

func (a *App) renderTemplate(w http.ResponseWriter, templateName string, data map[string]any) {
	// Construct the file path for the specified template
	tmplFile := filepath.Join("templates", templateName)

	//Read templates file into template set
	t, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data != nil {
		data["URL"] = a.Url
	}

	//execute template
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
