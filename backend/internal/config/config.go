// db config struct and loading
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Server struct {
		Port string
	}
	JWT struct {
		AccessSecret  string
		RefreshSecret string
	}
	OAuth struct {
		GoogleClientID     string
		GoogleClientSecret string
		GoogleRedirectURL  string
	}
}

func Load() (*Config, error) {
	// load .env file
	if err := godotenv.Load(filepath.Join("internal", "config", ".env")); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := &Config{}

	cfg.Database.Host = "db"
	cfg.Database.Port = "5432"
	cfg.Database.User = "user01"
	cfg.Database.Password = "user01239nTGN35pio!"
	cfg.Database.DBName = "reppydb"

	cfg.Server.Port = "8081"

	cfg.OAuth.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.OAuth.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	cfg.OAuth.GoogleRedirectURL = "http://localhost:8081/auth/google/callback"

	return cfg, nil
}
