package seeds

import (
	"context"
	"fmt"
	"go-reppy/backend/internal/database/sqlc"
)

// First add muscles, then the exercise_muscles to show which exercise_ids work which muscles, and whether primary or secondary.

type TestMuscles struct {
	MuscleName  string
	MuscleGroup string
}

func GetTestMuscles() []TestMuscles {
	return []TestMuscles{ // all major muscles/muscle groups
		// Chest
		{MuscleName: "Pectoralis Major", MuscleGroup: "Chest"},
		{MuscleName: "Pectoralis Minor", MuscleGroup: "Chest"},

		// Back
		{MuscleName: "Latissimus Dorsi", MuscleGroup: "Back"},
		{MuscleName: "Trapezius", MuscleGroup: "Back"},
		{MuscleName: "Rhomboids", MuscleGroup: "Back"},
		{MuscleName: "Teres Major", MuscleGroup: "Back"},
		{MuscleName: "Teres Minor", MuscleGroup: "Back"},

		// Shoulders
		{MuscleName: "Anterior Deltoid", MuscleGroup: "Shoulders"},
		{MuscleName: "Lateral Deltoid", MuscleGroup: "Shoulders"},
		{MuscleName: "Posterior Deltoid", MuscleGroup: "Shoulders"},

		// Arms
		{MuscleName: "Biceps Brachii", MuscleGroup: "Arms"},
		{MuscleName: "Triceps Brachii", MuscleGroup: "Arms"},
		{MuscleName: "Brachialis", MuscleGroup: "Arms"},
		{MuscleName: "Forearm Flexors", MuscleGroup: "Arms"},
		{MuscleName: "Forearm Extensors", MuscleGroup: "Arms"},

		// Legs
		{MuscleName: "Quadriceps", MuscleGroup: "Legs"},
		{MuscleName: "Hamstrings", MuscleGroup: "Legs"},
		{MuscleName: "Gastrocnemius", MuscleGroup: "Legs"},
		{MuscleName: "Soleus", MuscleGroup: "Legs"},
		{MuscleName: "Tibialis Anterior", MuscleGroup: "Legs"},

		// Core
		{MuscleName: "Rectus Abdominis", MuscleGroup: "Core"},
		{MuscleName: "Obliques", MuscleGroup: "Core"},
		{MuscleName: "Transverse Abdominis", MuscleGroup: "Core"},

		// Lower Back
		{MuscleName: "Erector Spinae", MuscleGroup: "Lower Back"},

		// Glutes
		{MuscleName: "Gluteus Maximus", MuscleGroup: "Glutes"},
		{MuscleName: "Gluteus Medius", MuscleGroup: "Glutes"},
		{MuscleName: "Gluteus Minimus", MuscleGroup: "Glutes"},
	}
}

func SeedMuscles(queries *sqlc.Queries) error {
	for _, muscle := range GetTestMuscles() {
		_, err := queries.CreateMuscle(context.Background(), sqlc.CreateMuscleParams{
			MuscleName:  muscle.MuscleName,
			MuscleGroup: muscle.MuscleGroup,
		})
		if err != nil {
			return fmt.Errorf("failed to create muscle %s: %v", muscle.MuscleName, err)
		}
	}
	fmt.Println("Successfully seeded MUSCLES table")
	return nil
}
