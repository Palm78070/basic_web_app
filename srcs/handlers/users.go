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

	// userList, err := app.Models.User.GetListWithUsername()
	// if err != nil {
	// 	fmt.Println("Error at GetListWithUsername")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// len_user := len(userList)
	// for i:=0; i < len_user; i++ {
	// 	tmp_user := userList[i]
	// 	fmt.Println(tmp_user.Id, tmp_user.Username.String, tmp_user.Email)
	// }

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

	// len_user := len(userList)
	// for i:=0; i < len_user; i++ {
	// 	tmp_user := userList[i]
	// 	fmt.Println(tmp_user.Id, tmp_user.Username.String, tmp_user.Email)
	// }

	app.renderTemplate(w, "userList.html", map[string]any{
		"UserList": userList,
	})
}
