package main

import (
	"log"
	"net/http"

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
		"scheme": "http://",
		"host": "localhost",
		"port": ":8080",
	}

	app := handlers.NewApp(config, db, url)

	log.Printf("Connected to db: %v", db)

	//Setting file server
	// Serve static files from the /static/ directory that locate on the root of the project
	//fs := http.FileServer(http.Dir("./static")) //Creates a handler to serve files from a specified directory
	//http.Handle("/static/", http.StripPrefix("/static/", fs)) //Registers the handler for a specific URL path StripPrefix => remove prefix static so result => static/css/file_to_serve
	// fs := http.FileServer(http.Dir("srcs/static"))
	// http.Handle("/static/", logRequest(setMimeType(http.StripPrefix("/static/", fs))))
	// http.Handle("/static/", fs)

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static")) //Creates a handler to serve files from a specified directory
	//http.Handle("/static/", http.StripPrefix("/static/", fs)) //Registers the handler for a specific URL path StripPrefix => remove prefix static so result => static/css/file_to_serve
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)) //match any path that has /static/ prefix handler is an http.Handler that will process the HTTP requests

	router.HandleFunc("/", app.IndexPage)
	router.HandleFunc("/register", app.RegisterPage).Methods("GET", "POST")
	router.HandleFunc("/login", app.Login)
	router.HandleFunc("/logout", app.Logout)
	router.HandleFunc("/callback", app.Callback)
	router.HandleFunc("/users/{username}", app.UserPage)
	router.HandleFunc("/userList", app.UserList)
	http.ListenAndServe(url["port"], router)
	// http.ListenAndServe(":8080", router)
}
