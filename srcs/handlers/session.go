package handlers

import (
	"fmt"

	"github.com/gorilla/sessions"
)

func (app *App)session_map_user(session *sessions.Session){
	// fmt.Println("In session_map_user")
	if session.Values["IsLogin"] != nil {
		app.currentUser.IsLogin = session.Values["IsLogin"].(bool)
	}
	if session.Values["username"] != nil {
		app.currentUser.username = session.Values["username"].(string)
	}
	if session.Values["email"] != nil {
		app.currentUser.email = session.Values["email"].(string)
	}
}

func (app *App)session_exist(session *sessions.Session) (bool){
	fmt.Println("In session_exist")
	_, ok := session.Values["IsLogin"]
	fmt.Println("session_exist: ", ok)
	return ok
}
