package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

func (app *App) registerUser(username string, email string, password string) {
	sqlStatement := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)"
	_, err := app.Models.User.DB.Exec(sqlStatement, username, email, password)
	if err != nil {
		fmt.Printf("Error inserting user into database: %v\n", err)
		return
	}
	fmt.Println("User inserted successfully")
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
		err := r.ParseMultipartForm(10 << 20) //Maximum form size 10MB
		if err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}


		user.Username = sql.NullString{
			String: r.FormValue("username"),
			Valid: true,
		}

		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")

		if !user.Username.Valid || user.Email == "" || user.Password == "" {
			log.Printf("Invalid form submission\n")
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"message": "Invalid form submission, username or email is not valid",
			})
			return
		}

		if app.is_register_user(user.Username.String, user.Email) {
			log.Printf("User already exists\n")
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"message": "User already exists, please use other username or email",
			})
			return
		}

		fmt.Println("username: ", user.Username)
		fmt.Println("password: ", user.Password)
		fmt.Println("email: ", user.Email)

		app.registerUser(user.Username.String, user.Email, user.Password)
		// http.Redirect(w, r, "/register", http.StatusSeeOther)
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
		})
	}
}
