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

	data["URL"] = a.Url

	//execute template
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//find all templates in templates folder
	// templates, err := filepath.Glob("templates/*.html")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	//Read all templates into template set
	// t, err := template.ParseFiles(templates...)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	//execute template
	// t.ExecuteTemplate(w, templateName, data)
}
