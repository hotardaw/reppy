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
	ResistanceType   *string // 'weight', 'band', 'bodyweight'
	ResistanceValue  *int32  // weight in pounds/kgs
	ResistanceDetail *string // "blue band", "wide grip", "olympic bar", etc.
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
			ResistanceValue:  utils.IntPtr(0),
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
			ResistanceValue:  utils.IntPtr(0),
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
			ResistanceValue:  utils.IntPtr(0),
			ResistanceDetail: utils.StrPtr("Wide grip"),
			RPE:              utils.Float32Ptr(8.5),
			Notes:            utils.StrPtr("Barely 7"),
		},
	}
}

// Failed to seed test data: failed to seed workouts: failed to seed workout-sets 1, 4, 1: pq: invalid input value for enum resistance_type_enum: "Bodyweight, +0 lbs"
func SeedWorkoutSets(queries *sqlc.Queries) error {
	for _, set := range GetTestWorkoutSets() {
		_, err := queries.CreateWorkoutSet(context.Background(), sqlc.CreateWorkoutSetParams{
			WorkoutID:        set.WorkoutID,
			ExerciseID:       set.ExerciseID,
			SetNumber:        set.SetNumber,
			Reps:             utils.NullIntFromIntPtr(set.Reps),
			ResistanceValue:  utils.NullIntFromIntPtr(set.ResistanceValue),
			ResistanceType:   utils.NullResistanceTypeEnumFromStringPtr(set.ResistanceType),
			ResistanceDetail: utils.NullStringFromStringPtr(set.ResistanceDetail),
			Rpe:              utils.NullStringFromFloat32Ptr(set.RPE),
			Notes:            utils.NullStringFromStringPtr(set.Notes),
		})
		if err != nil {
			return fmt.Errorf("failed to seed workout-sets %d, %d, %d: %v", set.WorkoutID, set.ExerciseID, set.SetNumber, err)
		}
	}
	fmt.Println("Successfully seeded WORKOUT-SETS table")
	return nil
}
