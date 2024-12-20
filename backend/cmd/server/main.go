// https://gowebexamples.com/
// https://pkg.go.dev/net/http
// https://http.dev/1.1
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-fitsync/backend/internal/api/handlers"
	"go-fitsync/backend/internal/api/middleware"
	"go-fitsync/backend/internal/config"
	"go-fitsync/backend/internal/database"
	"go-fitsync/backend/internal/database/seeder"
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

	queries := sqlc.New(db)

	// Account seeding
	if err := seeder.SeedTestData(queries); err != nil {
		// Printf instead of Fatal so app continues even if seeding fails
		log.Printf("Warning: Failed to seed test data: %v", err)
	}

	// Set JWT config, initialize auth middleware
	jwtConfig := middleware.JWTConfig{
		AccessSecret:    []byte(cfg.JWT.AccessSecret),
		RefreshSecret:   []byte(cfg.JWT.RefreshSecret),
		AccessDuration:  15 * time.Minute,
		RefreshDuration: 7 * 24 * time.Hour,
		Issuer:          "fitsync",
	}
	authMiddleware := middleware.NewAuthMiddleware(jwtConfig)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(queries, authMiddleware)
	userHandler := handlers.NewUserHandler(queries)
	userByIDHandler := handlers.NewUserByIDHandler(queries)
	userProfileHandler := handlers.NewUserProfileHandler(queries)
	userProfileByIDHandler := handlers.NewUserProfileByIDHandler(queries)
	muscleHandler := handlers.NewMuscleHandler(queries)

	mux := http.NewServeMux()

	// Auth routes (unprotected)
	mux.HandleFunc("/login/", middleware.LoggingMiddleware(authHandler.HandleLogin))
	mux.HandleFunc("/refresh/", middleware.LoggingMiddleware(authHandler.HandleRefresh))

	// User routes (protected)
	mux.HandleFunc("/users", middleware.LoggingMiddleware(userHandler.HandleUsers))         // GET(all), POST
	mux.HandleFunc("/users/", middleware.LoggingMiddleware(userByIDHandler.HandleUserByID)) // GET, PATCH, DELETE w/ ID

	mux.HandleFunc("/user-profiles", middleware.LoggingMiddleware(userProfileHandler.HandleUserProfiles))          // GET(all), POST
	mux.HandleFunc("/user-profiles/", middleware.LoggingMiddleware(userProfileByIDHandler.HandleUserProfilesByID)) // GET, PATCH, DELETE w/ ID

	mux.HandleFunc("/muscles", middleware.LoggingMiddleware(muscleHandler.HandleMuscles))

	// mux.HandleFunc("/workouts", middleware.LoggingMiddleware(workoutHandler.HandleWorkouts))     // GET, POST
	// mux.HandleFunc("/workouts/", middleware.LoggingMiddleware(workoutHandler.HandleWorkoutByID)) // GET, PATCH, DELETE w/ ID

	// Default/root handler
	mux.HandleFunc("/", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><head><title>FitSync API</title></head><body>Hello, World!</body></html>")
	}))

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Fatal(err)
	}
}
