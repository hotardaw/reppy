// db connection setup
package database

import (
	"database/sql"
	"fmt"

	"go-reppy/backend/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	fmt.Println("Opening database connection...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err = db.Ping(); err != nil { // set err to db.Ping result, see if nil
		return nil, err
	}

	fmt.Println("Database open at ", connStr)
	return db, nil
}
