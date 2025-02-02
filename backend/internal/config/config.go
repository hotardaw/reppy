// db config struct and loading
package config

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
	// later: load from env vars or config file
	cfg := &Config{}
	cfg.Database.Host = "db"
	cfg.Database.Port = "5432"
	cfg.Database.User = "user01"
	cfg.Database.Password = "user01239nTGN35pio!"
	cfg.Database.DBName = "reppydb"

	cfg.Server.Port = "8081"

	cfg.OAuth.GoogleClientID = "1091007547452-vuervm6jrk4o9d8rf1m814ttlkpn6r2b.apps.googleusercontent.com"
	cfg.OAuth.GoogleClientSecret = "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"
	cfg.OAuth.GoogleRedirectURL = "http://localhost:8081/auth/google/callback"

	return cfg, nil
}
