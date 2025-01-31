// https://gowebexamples.com/
// https://pkg.go.dev/net/http
// https://http.dev/1.1
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"go-fitstat/backend/internal/api/handlers"
	"go-fitstat/backend/internal/api/middleware"
	"go-fitstat/backend/internal/config"
	"go-fitstat/backend/internal/database"
	"go-fitstat/backend/internal/database/seeder"
	"go-fitstat/backend/internal/database/sqlc"
)

const (
	timeoutDuration   = 10 * time.Second
	maxBodySize       = 1024 * 1024 // 1mb
	requestsPerSecond = 5
	burstSize         = 10
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
		Issuer:          "fitstat",
	}

	timeoutMiddleware := middleware.TimeoutMiddleware(timeoutDuration)
	loggingMiddleware := middleware.LoggingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(jwtConfig)
	maxBodySizeMiddleware := middleware.MaxBodySizeMiddleware(maxBodySize)
	rateLimitMiddleware := middleware.RateLimitMiddleware(requestsPerSecond, burstSize)

	baseMiddleware := []func(http.HandlerFunc) http.HandlerFunc{
		timeoutMiddleware,
		loggingMiddleware,
		maxBodySizeMiddleware,
		rateLimitMiddleware,
		// add active-status-only for users in the db, since we perform soft deletes on users
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
	workoutSetHandler := handlers.NewWorkoutSetHandler(queries, jwtConfig.AccessSecret)
	workoutSetByIDHandler := handlers.NewWorkoutSetByIDHandler(queries, jwtConfig.AccessSecret)

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/signup", unprotected(authHandler.HandleSignup))
	mux.HandleFunc("/login", unprotected(authHandler.HandleLogin))
	mux.HandleFunc("/refresh", unprotected(authHandler.HandleRefresh))
	mux.HandleFunc("/logout", unprotected(authHandler.HandleLogout))

	// User routes
	mux.HandleFunc("/users", protected(userHandler.HandleUsers))                                                          // GET(all), POST
	mux.HandleFunc("/users/", protected(userByIDHandler.HandleUserByID))                                                  // GET, PATCH, DELETE
	mux.HandleFunc("/user-profiles", protected(userProfileHandler.HandleUserProfiles))                                    // GET(all), GET(active), POST
	mux.HandleFunc("/user-profiles/", protected(userProfileByIDHandler.HandleUserProfilesByID))                           // GET, PATCH, DELETE
	mux.HandleFunc("/muscles", protected(muscleHandler.HandleMuscles))                                                    // GET, POST, DELETE
	mux.HandleFunc("/exercises", protected(exerciseHandler.HandleExercises))                                              // GET(all), POST
	mux.HandleFunc("/exercises/", protected(exerciseByIDHandler.HandleExercisesByID))                                     // GET, PATCH, DELETE
	mux.HandleFunc("/workouts/{workout_id}/workout-sets", protected(workoutSetHandler.HandleWorkoutSets))                 // POST, GET(all), DELETE
	mux.HandleFunc("/workouts/{workout_id}/workout-sets/{set_id}", protected(workoutSetByIDHandler.HandleWorkoutSetByID)) // PATCH, DELETE
	mux.HandleFunc("/workouts", protected(workoutHandler.HandleWorkouts))                                                 // GET(all),POST
	mux.HandleFunc("/workouts/", protected(workoutByIDHandler.HandleWorkoutsByID))                                        // GET, PATCH, DELETE

	// Default/root handler
	mux.HandleFunc("/", unprotected(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html>
<head><title>FitStat API</title></head>
<body>
<u>FitStat API Routes</u>:
<li>/signup</li>
<li>/login</li>
<li>/refresh</li>
<li>/logout</li>
<li>/users</li>
<li>/users/{id}</li>
<li>/user-profiles</li>
<li>/user-profiles/{id}</li>
<li>/muscles</li>
<li>/exercises</li>
<li>/exercises/{id}</li>
<li>/workouts</li>
<li>/workouts/{id}</li>
<li>/workouts/{workout_id}/workout-sets</li>
<li>/workouts/{workout_id}/workout-sets/{set_id}</li>
</body>
</html>`)
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
