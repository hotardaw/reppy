package seeds

import (
	"context"
	"fmt"
	"go-fitstat/backend/internal/api/utils"
	"go-fitstat/backend/internal/database/sqlc"
)

type TestExercises struct {
	ExerciseName string
	Description  string
}

func GetTestExercises() []TestExercises {
	return []TestExercises{
		// Compound Exercises - Upper Body Push
		{
			ExerciseName: "Bench Press",
			Description:  "A compound exercise performed lying on a bench, pushing a barbell up from the chest to full arm extension.",
		},
		{
			ExerciseName: "Overhead Press",
			Description:  "A standing compound exercise pressing a barbell or dumbbells from shoulder level to overhead.",
		},
		{
			ExerciseName: "Push-up",
			Description:  "A bodyweight exercise performed face-down, pushing the body up from the ground with arms.",
		},

		// Compound Exercises - Upper Body Pull
		{
			ExerciseName: "Pull-up",
			Description:  "A bodyweight exercise pulling oneself up to a bar from a hanging position.",
		},
		{
			ExerciseName: "Barbell Row",
			Description:  "A bent-over pulling movement targeting the back muscles using a barbell.",
		},
		{
			ExerciseName: "Lat Pulldown",
			Description:  "A cable exercise pulling a bar down to the upper chest, targeting the latissimus dorsi.",
		},

		// Compound Exercises - Lower Body
		{
			ExerciseName: "Squat",
			Description:  "A fundamental lower body exercise performing a deep knee bend while keeping the torso upright.",
		},
		{
			ExerciseName: "Deadlift",
			Description:  "A compound exercise lifting a barbell from the ground while maintaining a neutral spine.",
		},
		{
			ExerciseName: "Romanian Deadlift",
			Description:  "A hip-hinge movement performed with straight legs, targeting the posterior chain.",
		},
		{
			ExerciseName: "Lunge",
			Description:  "A unilateral leg exercise stepping forward into a split stance position.",
		},

		// Isolation Exercises - Upper Body
		{
			ExerciseName: "Bicep Curl",
			Description:  "An isolation exercise for the biceps, curling weight from full arm extension to maximum flexion.",
		},
		{
			ExerciseName: "Tricep Extension",
			Description:  "An isolation movement extending the arm to target the triceps.",
		},
		{
			ExerciseName: "Lateral Raise",
			Description:  "An isolation exercise raising dumbbells to the side to target the lateral deltoids.",
		},

		// Isolation Exercises - Lower Body
		{
			ExerciseName: "Leg Extension",
			Description:  "A machine exercise extending the knee to target the quadriceps.",
		},
		{
			ExerciseName: "Leg Curl",
			Description:  "A machine exercise curling the leg to target the hamstrings.",
		},
		{
			ExerciseName: "Calf Raise",
			Description:  "An isolation exercise rising onto the toes to target the calf muscles.",
		},

		// Core Exercises
		{
			ExerciseName: "Plank",
			Description:  "An isometric core exercise maintaining a straight body position supported on forearms and toes.",
		},
		{
			ExerciseName: "Russian Twist",
			Description:  "A rotational core exercise performed seated with the feet off the ground.",
		},
		{
			ExerciseName: "Crunch",
			Description:  "A basic abdominal exercise lifting the shoulders off the ground while lying on the back.",
		},
	}
}

func SeedExercises(queries *sqlc.Queries) error {
	for _, exercise := range GetTestExercises() {
		_, err := queries.CreateExercise(context.Background(), sqlc.CreateExerciseParams{
			ExerciseName: exercise.ExerciseName,
			Description:  utils.ToNullString(exercise.Description),
		})
		if err != nil {
			return fmt.Errorf("failed to seed exercise %s: %v", exercise.ExerciseName, err)
		}
	}
	fmt.Println("Successfully seeded EXERCISES table")
	return nil
}
