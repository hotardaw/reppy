package seeds

import (
	"context"
	"fmt"
	"go-reppy/backend/internal/api/utils"
	"go-reppy/backend/internal/database/sqlc"
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

func SeedWorkouts(queries *sqlc.Queries) error {
	for _, workouts := range GetTestWorkouts() {
		_, err := queries.CreateWorkout(context.Background(), sqlc.CreateWorkoutParams{
			UserID:      utils.ToNullInt32(workouts.UserID),
			WorkoutDate: workouts.WorkoutDate,
			Title:       utils.ToNullString(workouts.Title),
		})
		if err != nil {
			return fmt.Errorf("failed to seed workout for UserID %d and WorkoutDate %v: %v", workouts.UserID, workouts.WorkoutDate, err)
		}
	}
	fmt.Println("Successfully seeded WORKOUTS table")
	return nil
}
