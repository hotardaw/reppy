package seeds

import (
	"context"
	"fmt"
	"go-fitsync/backend/internal/database/sqlc"
)

type TestExerciseMuscles struct {
	ExerciseID       int32
	MuscleID         int32
	InvolvementLevel string
}

func GetTestExerciseMuscles() []TestExerciseMuscles {
	return []TestExerciseMuscles{
		// Compound Exercises - Upper Body Push
		{
			ExerciseID:       1, // ID for bench press
			MuscleID:         1, // ID for pec major
			InvolvementLevel: "Primary",
		},
		{
			ExerciseID:       1, // ID for "Bench Press"
			MuscleID:         8, // ID for "Anterior Deltoid"
			InvolvementLevel: "Secondary",
		},
	}
}

func SeedExerciseMuscles(queries *sqlc.Queries) error {
	for _, exerciseMuscle := range GetTestExerciseMuscles() {
		_, err := queries.CreateExerciseMuscle(context.Background(), sqlc.CreateExerciseMuscleParams{
			ExerciseID:       exerciseMuscle.ExerciseID,
			MuscleID:         exerciseMuscle.MuscleID,
			InvolvementLevel: exerciseMuscle.InvolvementLevel,
		})
		if err != nil {
			return fmt.Errorf("failed to seed exercise-muscle association %d-%d: %v", exerciseMuscle.ExerciseID, exerciseMuscle.MuscleID, err)
		}
	}
	fmt.Println("Successfully seeded EXERCISE-MUSCLES table")
	return nil
}
