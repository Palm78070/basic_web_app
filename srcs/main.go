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

	app := handlers.NewApp(config, db)

	log.Printf("Connected to db: %v", db)

	router := mux.NewRouter()
	router.HandleFunc("/", app.IndexPage)
	router.HandleFunc("/users/{username}", app.UserPage)
	http.ListenAndServe(":8080", router)
}
