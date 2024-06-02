package handlers

import (
	"database/sql"

	"github.com/Palm78070/basic_web_app/models"
	"github.com/Palm78070/basic_web_app/settings"
)

type Models struct {
	User *models.UserModel
}

type App struct {
	Settings *settings.Settings
	Models *Models
}

func NewApp(config *settings.Settings, db *sql.DB) *App {
	return &App{
		Settings: config,
		Models: &Models{
			User: &models.UserModel{DB: db},
		},
	}
}
