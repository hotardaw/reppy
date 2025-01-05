package seeds

import (
	"context"
	"fmt"

	"go-fitsync/backend/internal/api/utils"
	"go-fitsync/backend/internal/database/sqlc"
)

type TestWorkoutSets struct { // Fields with pointers are nil-able
	WorkoutID        int32
	ExerciseID       int32
	SetNumber        int32 // Incrementing logic will need to be handled in app layer
	Reps             *int32
	ResistanceType   *string  // 'weight', 'band', 'bodyweight'
	ResistanceValue  *float32 // weight in pounds/kgs
	ResistanceDetail *string  // "blue band", "wide grip", "olympic bar", etc.
	RPE              *float32
	Notes            *string
}

func GetTestWorkoutSets() []TestWorkoutSets {
	return []TestWorkoutSets{
		{
			WorkoutID:  1,
			ExerciseID: 1,
			SetNumber:  1,
		},
		{
			WorkoutID:  1,
			ExerciseID: 1,
			SetNumber:  2,
		},
		{
			WorkoutID:  1,
			ExerciseID: 1,
			SetNumber:  3,
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        1,
			Reps:             utils.IntPtr(8),
			ResistanceType:   utils.StrPtr("bodyweight"),
			ResistanceValue:  utils.Float32Ptr(0),
			ResistanceDetail: utils.StrPtr("Wide grip"),
			RPE:              utils.Float32Ptr(8.5),
			Notes:            utils.StrPtr("Easy af"),
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        2,
			Reps:             utils.IntPtr(8),
			ResistanceType:   utils.StrPtr("bodyweight"),
			ResistanceValue:  utils.Float32Ptr(0),
			ResistanceDetail: utils.StrPtr("Wide grip"),
			RPE:              utils.Float32Ptr(8.5),
			Notes:            utils.StrPtr("Less easy this time"),
		},
		{
			WorkoutID:        1,
			ExerciseID:       4,
			SetNumber:        3,
			Reps:             utils.IntPtr(7),
			ResistanceType:   utils.StrPtr("bodyweight"),
			ResistanceValue:  utils.Float32Ptr(0),
			ResistanceDetail: utils.StrPtr("Wide grip"),
			RPE:              utils.Float32Ptr(8.5),
			Notes:            utils.StrPtr("Barely 7"),
		},
	}
}

// Failed to seed test data: failed to seed workouts: failed to seed workout-sets 1, 4, 1: pq: invalid input value for enum resistance_type_enum: "Bodyweight, +0 lbs"
func SeedWorkoutSets(queries *sqlc.Queries) error {
	sets := GetTestWorkoutSets()

	workoutIDs := make([]int32, len(sets))
	exerciseIDs := make([]int32, len(sets))
	setNumbers := make([]int32, len(sets))
	reps := make([]int32, len(sets))
	resistanceValues := make([]string, len(sets))
	resistanceTypes := make([]string, len(sets))
	resistanceDetails := make([]string, len(sets))
	rpes := make([]string, len(sets)) // expecting numeric input
	notes := make([]string, len(sets))

	for i, set := range sets {
		workoutIDs[i] = set.WorkoutID
		exerciseIDs[i] = set.ExerciseID
		setNumbers[i] = set.SetNumber
		if set.Reps != nil {
			reps[i] = *set.Reps
		}
		if set.ResistanceValue != nil {
			resistanceValues[i] = fmt.Sprintf("%.1f", *set.ResistanceValue)
		}
		if set.ResistanceType != nil {
			resistanceTypes[i] = *set.ResistanceType
		}
		if set.ResistanceDetail != nil {
			resistanceDetails[i] = *set.ResistanceDetail
		}
		if set.RPE != nil {
			rpes[i] = fmt.Sprintf("%.1f", *set.RPE)
		}
		if set.Notes != nil {
			notes[i] = *set.Notes
		}
	}

	_, err := queries.CreateWorkoutSets(context.Background(), sqlc.CreateWorkoutSetsParams{
		Column1: workoutIDs,
		Column2: exerciseIDs,
		Column3: setNumbers,
		Column4: reps,
		Column5: resistanceValues,
		Column6: resistanceTypes,
		Column7: resistanceDetails,
		Column8: rpes,
		Column9: notes,
	})
	if err != nil {
		return fmt.Errorf("failed to seed workout-sets: %v", err)
	}

	fmt.Println("Successfully seeded WORKOUT-SETS table")
	return nil
}
