package main

import (
	"log"
	"net/http"

	"github.com/Palm78070/basic_web_app/db"
	"github.com/Palm78070/basic_web_app/handlers"
	"github.com/Palm78070/basic_web_app/settings"

	"github.com/gorilla/mux"
)

func main() {
	config, err := settings.LoadSettings()

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

	router := mux.NewRouter()
	router.HandleFunc("/", app.IndexPage)
	router.HandleFunc("/users/{username}", app.UserPage)
	router.HandleFunc("/userList", app.UserList)
	http.ListenAndServe(url["port"], router)
	// http.ListenAndServe(":8080", router)
}
