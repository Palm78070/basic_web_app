package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Palm78070/basic_web_app/models"
)

func (app *App) is_register_user(username string, email string) bool {
	var result bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)"
	err := app.Models.User.DB.QueryRow(query, username, email).Scan(&result)
	if err != nil && err != sql.ErrNoRows{
		fmt.Printf("error checking user existence: %v\n", err)
		return false
	}
	return result
}

func (app *App) registerUser(username string, email string) {
	sqlStatement := "INSERT INTO users (username, email) VALUES ($1, $2)"
	_, err := app.Models.User.DB.Exec(sqlStatement, username, email)
	if err != nil {
		fmt.Printf("Error inserting user into database: %v\n", err)
		return
	}
}

func (app *App) RegisterPage(w http.ResponseWriter, r* http.Request) {
	if r.Method == http.MethodGet {
		app.renderTemplate(w, "register.html", map[string]any{
			"URL": app.Url,
		})
		return
	}

	if r.Method == http.MethodPost {
		var user models.User
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}

		user.Username = sql.NullString{
			String: r.FormValue("username"),
			Valid: true,
		}

		user.Email = r.FormValue("email")

		if !user.Username.Valid || user.Email == "" {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}


		if app.is_register_user(user.Username.String, user.Email) {
			http.Error(w, "User already exists", http.StatusBadRequest)
		}

		fmt.Println("username: ", user.Username)
		fmt.Println("email: ", user.Email)

		app.registerUser(user.Username.String, user.Email)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
}
