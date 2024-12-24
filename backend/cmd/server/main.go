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

	if err := seeder.SeedTestData(queries); err != nil {
		// Printf instead of Fatal so app continues even if seeding fails
		log.Printf("Warning: Failed to seed test data: %v", err)
	}

	jwtConfig := middleware.JWTConfig{
		AccessSecret:    []byte(cfg.JWT.AccessSecret),
		RefreshSecret:   []byte(cfg.JWT.RefreshSecret),
		AccessDuration:  15 * time.Minute,
		RefreshDuration: 7 * 24 * time.Hour,
		Issuer:          "fitsync",
	}
	timeoutMiddleware := middleware.TimeoutMiddleware(10 * time.Second)
	loggingMiddleware := middleware.LoggingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(jwtConfig)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(queries, authMiddleware)
	userHandler := handlers.NewUserHandler(queries)
	userByIDHandler := handlers.NewUserByIDHandler(queries)
	userProfileHandler := handlers.NewUserProfileHandler(queries)
	userProfileByIDHandler := handlers.NewUserProfileByIDHandler(queries)
	muscleHandler := handlers.NewMuscleHandler(queries)
	exerciseHandler := handlers.NewExerciseHandler(queries)

	mux := http.NewServeMux()

	// Auth routes (unprotected)
	mux.HandleFunc("/login/", chainMiddleware(
		authHandler.HandleLogin,
		timeoutMiddleware,
		loggingMiddleware,
	))
	mux.HandleFunc("/refresh/", chainMiddleware(
		authHandler.HandleRefresh,
		timeoutMiddleware,
		loggingMiddleware,
	))

	// User routes (protected)
	mux.HandleFunc("/users", chainMiddleware( // GET(all), POST
		userHandler.HandleUsers,
		timeoutMiddleware,
		loggingMiddleware,
		authMiddleware.AuthenticateJWT,
	))
	mux.HandleFunc("/users/", chainMiddleware( // GET, PATCH, DELETE w/ ID
		userByIDHandler.HandleUserByID,
		timeoutMiddleware,
		loggingMiddleware,
		authMiddleware.AuthenticateJWT,
	))

	mux.HandleFunc("/user-profiles", chainMiddleware( // GET(all), POST
		userProfileHandler.HandleUserProfiles,
		timeoutMiddleware,
		loggingMiddleware,
		authMiddleware.AuthenticateJWT,
	))
	mux.HandleFunc("/user-profiles", chainMiddleware( // GET, PATCH, DELETE w/ ID
		userProfileByIDHandler.HandleUserProfilesByID,
		timeoutMiddleware,
		loggingMiddleware,
		authMiddleware.AuthenticateJWT,
	))

	mux.HandleFunc("/muscles", chainMiddleware(
		muscleHandler.HandleMuscles,
		timeoutMiddleware,
		loggingMiddleware,
		authMiddleware.AuthenticateJWT,
	))

	mux.HandleFunc("/exercises", chainMiddleware(
		exerciseHandler.HandleExercises,
		timeoutMiddleware,
		loggingMiddleware,
		authMiddleware.AuthenticateJWT,
	))

	// mux.HandleFunc("/workouts", middleware.LoggingMiddleware(workoutHandler.HandleWorkouts))     // GET, POST
	// mux.HandleFunc("/workouts/", middleware.LoggingMiddleware(workoutHandler.HandleWorkoutByID)) // GET, PATCH, DELETE w/ ID

	// Default/root handler
	mux.HandleFunc("/", chainMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><head><title>FitSync API</title></head><body>
        <u>FitSync API Routes</u>:
        <li>/login/</li>
        <li>/refresh/</li>
        <li>/users</li>
        <li>/users/{user_id}</li>
        <li>/user-profiles</li>
        <li>/user-profiles/{user_id}</li>
        <li>/muscles</li>
        <li>/exercises</li>
        <h2>in progress:</h2>
        <li>/workouts</li>
        <li>/workouts/{workout_id}</li>
        
        </body></html>`)
		},
		timeoutMiddleware,
		loggingMiddleware,
	))

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Fatal(err)
	}
}

func chainMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
