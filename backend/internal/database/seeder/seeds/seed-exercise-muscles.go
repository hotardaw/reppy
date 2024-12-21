package seeds

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

// func SeedExerciseMuscles(queries *sqlc.Queries) error {
// 	for _, exercise := range GetTestExerciseMuscles() {
// 		_, err := queries.CreateExercise(context.Background(), sqlc.CreateExerciseParams{
// 			ExerciseName: exercise.ExerciseName,
// 			Description: sql.NullString{
// 				String: exercise.Description,
// 				Valid:  true,
// 			},
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to seed exercise %s: %v", exercise.ExerciseName, err)
// 		}
// 	}
// 	fmt.Println("Successfully seeded EXERCISES table")
// 	return nil
// }
