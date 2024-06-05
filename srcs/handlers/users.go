package handlers

import (
	"fmt"
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

	app.renderTemplate(w, "userPage.html", map[string]any{
		"User": user,
	})
}

func (app *App) UserList(w http.ResponseWriter, r *http.Request) {
	userList, err := app.Models.User.GetListWithUsername()
	if err != nil {
		fmt.Println("Error at GetListWithUsername")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.renderTemplate(w, "userList.html", map[string]any{
		"UserList": userList,
	})
}
