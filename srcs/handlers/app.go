package handlers

import (
	"database/sql"
	"os"

	"github.com/Palm78070/basic_web_app/models"
	"github.com/Palm78070/basic_web_app/settings"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Models struct {
	User *models.UserModel
}

type Login struct {
	googleOauthConfig *oauth2.Config
	randomState string
	IsLogin bool
	content map[string]interface{}
	email string
}
type App struct {
	Settings *settings.Settings
	Models *Models
	Url string
	currentUser *Login
}

func NewApp(config *settings.Settings, db *sql.DB, url map[string]string) *App {
	return &App{
		Settings: config,
		Models: &Models{
			User: &models.UserModel{DB: db},
		},
		Url: url["scheme"] + url["host"] + url["port"],
		currentUser: &Login{
			googleOauthConfig: &oauth2.Config{
				RedirectURL: url["scheme"] + url["host"] + url["port"] + "/callback",
				ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
				Endpoint: google.Endpoint,
			},
			//TODO: randomize it
			randomState : "random",
			IsLogin: false,
		},
	}
}
