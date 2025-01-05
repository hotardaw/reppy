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

const (
	timeoutDuration = 10 * time.Second
	maxBodySize     = 1024 * 1024 // 1mb
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
		AccessDuration:  30 * time.Minute, // change to 15 later
		RefreshDuration: 7 * 24 * time.Hour,
		Issuer:          "fitsync",
	}

	timeoutMiddleware := middleware.TimeoutMiddleware(timeoutDuration)
	loggingMiddleware := middleware.LoggingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(jwtConfig)
	maxBodySizeMiddleware := middleware.MaxBodySizeMiddleware(maxBodySize)

	baseMiddleware := []func(http.HandlerFunc) http.HandlerFunc{
		timeoutMiddleware,
		loggingMiddleware,
		maxBodySizeMiddleware,
		// add rate limiter, ESPECIALLY for refresh tokens
		// add active-status-only for users in the db, since we perform soft deletes
	}
	protectedMiddleware := append([]func(http.HandlerFunc) http.HandlerFunc{
		authMiddleware.AuthenticateJWT,
	}, baseMiddleware...)

	protected := func(handler http.HandlerFunc) http.HandlerFunc {
		return chainMiddleware(handler, protectedMiddleware...)
	}

	unprotected := func(handler http.HandlerFunc) http.HandlerFunc {
		return chainMiddleware(handler, baseMiddleware...)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(queries, authMiddleware)
	userHandler := handlers.NewUserHandler(queries)
	userByIDHandler := handlers.NewUserByIDHandler(queries)
	userProfileHandler := handlers.NewUserProfileHandler(queries)
	userProfileByIDHandler := handlers.NewUserProfileByIDHandler(queries)
	muscleHandler := handlers.NewMuscleHandler(queries)
	exerciseHandler := handlers.NewExerciseHandler(queries)
	exerciseByIDHandler := handlers.NewExerciseByIDHandler(queries)
	workoutHandler := handlers.NewWorkoutHandler(queries, jwtConfig.AccessSecret)
	workoutByIDHandler := handlers.NewWorkoutByIDHandler(queries)

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/signup", unprotected(authHandler.HandleSignup))
	mux.HandleFunc("/login", unprotected(authHandler.HandleLogin))
	mux.HandleFunc("/refresh", unprotected(authHandler.HandleRefresh))
	mux.HandleFunc("/logout", unprotected(authHandler.HandleLogout))

	// User routes
	mux.HandleFunc("/users", protected(userHandler.HandleUsers))                                // GET(all), POST
	mux.HandleFunc("/users/", protected(userByIDHandler.HandleUserByID))                        // GET, PATCH, DELETE w/ ID
	mux.HandleFunc("/user-profiles", protected(userProfileHandler.HandleUserProfiles))          // GET(all), GET(active), POST
	mux.HandleFunc("/user-profiles/", protected(userProfileByIDHandler.HandleUserProfilesByID)) // GET, PATCH, DELETE w/ ID
	mux.HandleFunc("/muscles", protected(muscleHandler.HandleMuscles))                          // GET, POST, DELETE
	mux.HandleFunc("/exercises", protected(exerciseHandler.HandleExercises))                    // GET(all), POST
	mux.HandleFunc("/exercises/", protected(exerciseByIDHandler.HandleExercisesByID))           // GET, PATCH, DELETE
	mux.HandleFunc("/workouts", protected(workoutHandler.HandleWorkouts))                       // GET(all),POST
	mux.HandleFunc("/workouts/", protected(workoutByIDHandler.HandleWorkoutsByID))              // GET, PATCH, DELETE w/ ID

	// Default/root handler
	mux.HandleFunc("/", unprotected(func(w http.ResponseWriter, r *http.Request) {
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
        <li>/workouts</li>
        <li>/workouts/{workout_id}</li>
				
        <h2>in progress:</h2>
        <li>/workout-sets</li>
        <li>/workout-sets/{workout_id}</li>
				
        </body></html>`)
	}))

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
