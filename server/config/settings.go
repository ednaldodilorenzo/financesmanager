package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	ENVIRONMENT_DEV  = "DEVELOPMENT"
	ENVIRONMENT_PROD = "PRODUCTION"
)

type DBSettings struct {
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbname"`
	LogLevel string
}

type BrokerSettings struct {
	Host string
	Port string
}

type AppSettings struct {
	Url         string
	JwtKey      string
	Environment string
}

type Settings struct {
	Database            DBSettings `yaml:"database"`
	MessageBroker       BrokerSettings
	AppSettings         AppSettings
	GoogleOAuthSettings oauth2.Config
}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) LoadSettings() {
	s.Database = DBSettings{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}

	s.MessageBroker = BrokerSettings{
		Host: os.Getenv("MB_HOST"),
		Port: os.Getenv("MB_PORT"),
	}

	s.AppSettings = AppSettings{
		Url:         os.Getenv("APP_URL"),
		JwtKey:      os.Getenv("APP_JWT_KEY"),
		Environment: os.Getenv("APP_ENVIRONMENT"),
	}

	s.GoogleOAuthSettings = oauth2.Config{
		RedirectURL:  "http://localhost:5000/api/auth/google/callback",
		ClientID:     "864370761747-8ranshuk43vcu5pm69uk0q0na40s77rg.apps.googleusercontent.com", //os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: "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",                                      //os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
