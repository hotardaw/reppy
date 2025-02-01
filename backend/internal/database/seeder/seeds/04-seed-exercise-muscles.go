package seeds

import (
	"context"
	"fmt"
	"go-reppy/backend/internal/database/sqlc"
)

type TestExerciseMuscles struct {
	ExerciseID       int32
	MuscleID         int32
	InvolvementLevel string
}

func GetTestExerciseMuscles() []TestExerciseMuscles {
	return []TestExerciseMuscles{
		// Bench Press
		{
			ExerciseID:       1,
			MuscleID:         1, // ID for "Pec Major"
			InvolvementLevel: "primary",
		},
		{
			ExerciseID:       1,
			MuscleID:         8, // ID for "Anterior Deltoid"
			InvolvementLevel: "secondary",
		},
		{
			ExerciseID:       1,
			MuscleID:         11, // Triceps
			InvolvementLevel: "secondary",
		},

		// Overhead Press
		{
			ExerciseID:       2,
			MuscleID:         8, // Anterior Deltoid
			InvolvementLevel: "primary",
		},
		{
			ExerciseID:       2,
			MuscleID:         11, // Triceps
			InvolvementLevel: "secondary",
		},

		// Push-up
		{
			ExerciseID:       3,
			MuscleID:         1, // Pectoralis Major
			InvolvementLevel: "primary",
		},
		{
			ExerciseID:       3,
			MuscleID:         11, // Triceps
			InvolvementLevel: "secondary",
		},
		{
			ExerciseID:       3,
			MuscleID:         8, // Anterior Deltoid
			InvolvementLevel: "secondary",
		},

		// Pull-up
		{
			ExerciseID:       4,
			MuscleID:         6, // Latissimus Dorsi
			InvolvementLevel: "primary",
		},
		{
			ExerciseID:       4,
			MuscleID:         12, // Biceps
			InvolvementLevel: "secondary",
		},
		{
			ExerciseID:       4,
			MuscleID:         7, // Posterior Deltoid
			InvolvementLevel: "secondary",
		},

		// Barbell Row
		{
			ExerciseID:       5,
			MuscleID:         6, // Latissimus Dorsi
			InvolvementLevel: "primary",
		},
		{
			ExerciseID:       5,
			MuscleID:         12, // Biceps
			InvolvementLevel: "secondary",
		},
		{
			ExerciseID:       5,
			MuscleID:         7, // Posterior Deltoid
			InvolvementLevel: "secondary",
		},
	}
}

func SeedExerciseMuscles(queries *sqlc.Queries) error {
	for _, exerciseMuscle := range GetTestExerciseMuscles() {
		_, err := queries.CreateExerciseMuscle(context.Background(), sqlc.CreateExerciseMuscleParams{
			ExerciseID:       exerciseMuscle.ExerciseID,
			MuscleID:         exerciseMuscle.MuscleID,
			InvolvementLevel: sqlc.InvolvementLevelEnum(exerciseMuscle.InvolvementLevel),
		})
		if err != nil {
			return fmt.Errorf("failed to seed exercise-muscle association %d-%d: %v", exerciseMuscle.ExerciseID, exerciseMuscle.MuscleID, err)
		}
	}
	fmt.Println("Successfully seeded EXERCISE-MUSCLES table")
	return nil
}
