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

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
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
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

