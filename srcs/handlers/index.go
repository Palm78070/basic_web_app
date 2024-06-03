package handlers

import (
	"net/http"
)

func (app *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "index.html", map[string]any{
	})
}
