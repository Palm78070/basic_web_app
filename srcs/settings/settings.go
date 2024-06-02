package settings

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type DBSettings struct {
	Host string `koanf:"HOST" validate:"required"` //DB__HOST
	Port int `koanf:"PORT" validate:"required"`
	User string `koanf:"USER" validate:"required"`
	Password string `koanf:"PASSWORD" validate:"required"`
	DBName string `koanf:"DBNAME" validate:"required"`
}

type Settings struct {
	Host string `koanf:"HOST" validate:"required"`
	Port int `koanf:"PORT" validate:"required"`
	DB DBSettings `koanf:"DB" validate:"required"`
}

func LoadSettings() (settings *Settings, err error) {
	k := koanf.New(".") //Use "." as the config path delimiter
	//Load setting from .env file
	_ = k.Load(
		file.Provider(".env"),
		dotenv.ParserEnv("", "__", func(s string) string {
			return s //Change to anonymous func for flexibility of customization in future
		}),
	)

	// "__" is for nesting
	//Load setting from .env variables
	_ = k.Load(env.Provider("", "__", nil), nil)

	// var settings Settings
	//Write data which load from .env file or env var to Settings struct
	err = k.Unmarshal("", &settings)
	if err != nil {
		return nil, err
	}

	//TODO: validate settings
	validate := validator.New()
	err = validate.Struct(settings) //validate a struct Settings and also nested struct DBSettings
	if err != nil {
		return nil, err
	}

	return settings, nil
}
