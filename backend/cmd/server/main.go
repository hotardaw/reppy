// https://gowebexamples.com/
// https://pkg.go.dev/net/http
// https://http.dev/1.1
package main

import (
	"fmt"
	"log"
	"net/http"

	"go-fitsync/backend/internal/api/handlers"
	"go-fitsync/backend/internal/config"
	"go-fitsync/backend/internal/database"
	"go-fitsync/backend/internal/database/sqlc"

	// this is "side-effect import".
	// it registers the pq (postgresql) driver with Go's database/sql package w/o directly using its exported identifiers
	// when imported, its init() runs, registering the driver with the database/sql package, allowing postgresql use with the standard database/sql interface w/o direct calls from the pq package.
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize queries
	queries := sqlc.New(db)

	mux := http.NewServeMux()

	// authHandler := handlers.NewAuthHandler(queries)
	// // Auth routes
	// mux.HandleFunc("/login/", authHandler.HandleLogin)
	// mux.HandleFunc("/refresh/", authHandler.HandleRefresh) // refresh token

	userHandler := handlers.NewUserHandler(queries)
	userProfileHandler := handlers.NewUserProfileHandler(queries)

	// User routes (protected)
	mux.HandleFunc("/users/", userHandler.HandleUsers)
	mux.HandleFunc("/user-profiles/", userProfileHandler.HandleUserProfiles)

	// Default/root handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><head><title>API</title></head><body>Hello, World!</body></html>")
	})

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Fatal(err)
	}
}
