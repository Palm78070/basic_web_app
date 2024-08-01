package handlers

import (
	"fmt"
	"net/http"
)

func (app *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("In IndexPage")
	session, _ := app.SessionStore.Get(r, "session-name")
	if !app.session_exist(session) {
		// http.Redirect(w, r, "/login", http.StatusForbidden)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.session_map_user(session)
	username := ""
	if app.currentUser.IsLogin {
		email := app.currentUser.email
		user, err := app.Models.User.GetByEmail(email)
		if err != nil {
			fmt.Printf("Cannot get user by email: %s\n", err.Error())
			http.Error(w, "Cannot get user by email", http.StatusInternalServerError)
			return
		}
		username = user.Username.String
		// fmt.Println("Email: ", email)
		// fmt.Println(user)
	}
	// fmt.Println("After get username: ", username)
	// fmt.Println("isLogin: ", app.currentUser.IsLogin)
	app.renderTemplate(w, "index.html", map[string]any{
		"IsLogin": app.currentUser.IsLogin,
		"Username": username,
	})
}
