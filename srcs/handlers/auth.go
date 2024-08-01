package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func (app *App)Login(w http.ResponseWriter, r *http.Request) {
	//session-name to retrieve session if exist or create a new one
	session, _ := app.SessionStore.Get(r, "session-name")
	if r.Method == http.MethodGet {
		app.renderTemplate(w, "login.html", map[string]any{
		})
		return
	}
	if r.Method != http.MethodPost {
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("Username: ", username)
	fmt.Println("Password: ", password)

	//row := m.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", username)
	ok, err := app.Models.User.CheckAuth(username, password)
	if !ok && err == nil {
		fmt.Println("Invalid username or password")
		// http.Error(w, "Invalid username or password", http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}
	if !ok && err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// app.currentUser.username = username
	// app.currentUser.IsLogin = true
	// app.currentUser.email = app.Models.User.GetEmail(username)
	// if app.currentUser.email == "" {
	// 	fmt.Println("Cannot get email")
	// 	http.Error(w, "Cannot get email", http.StatusInternalServerError)
	// 	return
	// }

	session.Values["username"] = username
	session.Values["IsLogin"] = true
	session.Values["email"] = app.Models.User.GetEmail(username)
	if session.Values["email"] == "" {
		fmt.Println("Cannot get email")
		http.Error(w, "Cannot get email", http.StatusInternalServerError)
		return
	}
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Cannot save session when login: ", err)
		http.Error(w, "Cannot save session when logging in", http.StatusInternalServerError)
		return
	}
	fmt.Println("IsLogin value in session: ", session.Values["IsLogin"])
	app.session_map_user(session)
	fmt.Println("User", username, "logged in successfully")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App)LoginGoogle(w http.ResponseWriter, r *http.Request) {
	//Redirect to google specific url
	// url := app.googleOauthConfig.AuthCodeURL(app.randomState)
	//prompt=login=>forces Google to display the login screen, even if the user is already logged in with their Google account

	url := app.currentUser.googleOauthConfig.AuthCodeURL(app.currentUser.randomState, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "login"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *App)Logout(w http.ResponseWriter, r *http.Request) {
	for key := range app.currentUser.content {
		delete(app.currentUser.content, key)
	}
	app.currentUser.IsLogin = false

	//Clear session data on server side
	session, _ := app.SessionStore.Get(r, "session-name")
	for k := range session.Values {
		delete(session.Values, k)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session-name",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})


	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (app *App)Callback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != app.currentUser.randomState {
		fmt.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//Provide context and code and retrieve token
	//https://www.notion.so/Golang-48669267a36e4db1863df9bbcb711716?pvs=4#76fb414760cb44d8a61a92eb811b431e
	token, err := app.currentUser.googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create get request: %s\n", err.Error())
		return
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not parse response: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// fmt.Fprintf(w, "Response: %s", content)
	// Redirect to the home page ("/").
	contentString := string(content)
	err = json.Unmarshal(content, &app.currentUser.content)
	if err != nil {
		fmt.Printf("could not parse JSON: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// app.currentUser.email
	emailInterface := app.currentUser.content["email"]
	email, ok := emailInterface.(string)
	if !ok {
		log.Println("email is not a string")
		http.Error(w, "Cannot get user by email, email is not a string", http.StatusInternalServerError)
		return
	}
	app.currentUser.email = email

	fmt.Println("Response: ", contentString)
	if contentString != "" {
		app.currentUser.IsLogin = true
	}
	user, err := app.Models.User.GetByEmail(app.currentUser.email)
	if err != nil {
		fmt.Println(err)
		return
	}
	if user == nil {
		fmt.Println("User not found")
		return
	}
	app.currentUser.username = user.Username.String
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

