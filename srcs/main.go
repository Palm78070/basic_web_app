package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Palm78070/basic_web_app/db"
	"github.com/Palm78070/basic_web_app/handlers"
	"github.com/Palm78070/basic_web_app/settings"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	sig_ch := make(chan os.Signal, 1)
	// signal.Notify(sig_ch, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sig_ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	done := make(chan struct{})

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
	go app.HandleMessages()
	go app.Cleanup(sig_ch, done)
	log.Fatal(http.ListenAndServe(url["port"], router))

	<-done
	log.Println("Server stopped")
	close(sig_ch)
	os.Exit(0)
}
