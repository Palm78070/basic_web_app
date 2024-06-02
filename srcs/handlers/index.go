package handlers

import (
	"net/http"
)

func (app *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}
