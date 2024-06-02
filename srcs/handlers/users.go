package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *App) UserPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := app.Models.User.GetByUsername(vars["username"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil || !user.Username.Valid {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "userPage.html", map[string]any{
		"User": user,
	})
}
