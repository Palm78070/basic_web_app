package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Palm78070/basic_web_app/db"
	"github.com/Palm78070/basic_web_app/handlers"
	"github.com/Palm78070/basic_web_app/settings"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	config, err := settings.LoadSettings()
	if err != nil {
		log.Fatalf("Error loading settings: %v", err)
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}

	db, err := db.Connect(&config.DB)
	if err != nil {
		panic(err)
	}

	url := map[string]string{
		"scheme": os.Getenv("SCHEME"),
		"host": os.Getenv("HOST"),
		"port": ":" + os.Getenv("PORT"),
	}

	app := handlers.NewApp(config, db, url)

	log.Printf("Connected to db: %v", db)

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static")) //Creates a handler to serve files from a specified directory
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)) //match any path that has /static/ prefix handler is an http.Handler that will process the HTTP requests

	router.HandleFunc("/", app.IndexPage)
	router.HandleFunc("/register", app.RegisterPage).Methods("GET", "POST")
	router.HandleFunc("/login", app.Login).Methods("GET", "POST")
	router.HandleFunc("/login/google", app.LoginGoogle)
	router.HandleFunc("/logout", app.Logout)
	router.HandleFunc("/callback", app.Callback)
	router.HandleFunc("/users/{username}", app.UserPage)
	router.HandleFunc("/userList", app.UserList)
	router.HandleFunc("/ws/{room_type}/{room_name}", app.HandleConnections)
	// go app.HandleMessages()
	http.ListenAndServe(url["port"], router)
	// http.ListenAndServe(":8080", router)
}
