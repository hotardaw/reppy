package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"go-fitsync/backend/internal/database/sqlc"
	"time"
)

type TestWorkouts struct {
	UserID      int32
	WorkoutDate time.Time
	Title       string
}

func GetTestWorkouts() []TestWorkouts {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	return []TestWorkouts{
		{
			UserID:      1,
			WorkoutDate: today,
			Title:       "Upper Body 1",
		},
		{
			UserID:      1,
			WorkoutDate: today.AddDate(0, 0, -1), // yesterday
			Title:       "Lower Body 1",
		},
	}
}

/*
-- CREATE: Insert a new workout
-- name: CreateWorkout :one

-- READ: Get all workouts for a specific user
-- name: GetAllWorkoutsForUser :many

-- READ: Get a specific workout by ID
-- name: GetWorkoutByID :one

-- READ: Get workouts within a date range for a user
-- name: GetWorkoutsWithinDateRange :many

-- UPDATE: Modify an existing workout
-- name: UpdateWorkout :one

-- DELETE: Remove a workout
-- name: DeleteWorkout :one

*/

func SeedWorkouts(queries *sqlc.Queries) error {
	for _, workouts := range GetTestWorkouts() {
		_, err := queries.CreateWorkout(context.Background(), sqlc.CreateWorkoutParams{
			UserID: sql.NullInt32{
				Int32: workouts.UserID,
				Valid: true,
			},
			WorkoutDate: workouts.WorkoutDate,
			Title: sql.NullString{
				String: workouts.Title,
				Valid:  true,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to seed workout for UserID %d and WorkoutDate %v: %v", workouts.UserID, workouts.WorkoutDate, err)
		}
	}
	fmt.Println("Successfully seeded WORKOUTS table")
	return nil
}
