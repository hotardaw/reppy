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
}

func Load() (*Config, error) {
	// later: load from env vars or config file
	cfg := &Config{}
	cfg.Database.Host = "db"
	cfg.Database.Port = "5432"
	cfg.Database.User = "user01"
	cfg.Database.Password = "user01239nTGN35pio!$"
	cfg.Database.DBName = "fitsyncdb"
	cfg.Server.Port = "8081"
	return cfg, nil
}
